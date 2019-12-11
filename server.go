package possum

import (
	"container/list"
	"net/http"
	"net/url"
	"sync"

	"github.com/mikespook/possum/router"
	"github.com/mikespook/possum/view"
)

type HandlerFunc func(w http.ResponseWriter, req *http.Request) (interface{}, int)

// Possum implements `http.Handler`.
type Possum struct {
	sync.RWMutex
	direct map[string]*rvcPack
	other  *list.List

	// PreRequest is called after initialised the request and session, before user-defined handler
	PreRequest http.HandlerFunc
	// PreResponse is called before sending response to client.
	PreResponse http.HandlerFunc
	// ErrorHandler gets a chance to write user-defined error handler
	ErrorHandle func(error)
}

// rvcPack shorts for router/view/controller pack
type rvcPack struct {
	router router.Router
	view   view.View
	f      HandlerFunc
}

// New initailizes Possum instance.
func New() *Possum {
	return &Possum{
		direct: make(map[string]*rvcPack),
		other:  list.New(),
	}
}

// Add a router to list
func (psm *Possum) Add(rtr router.Router, f HandlerFunc, view view.View) {
	defer psm.Unlock()
	psm.Lock()
	pack := &rvcPack{rtr, view, f}
	// direct will full-match the path
	if baseRouter, ok := rtr.(*router.Base); ok {
		psm.direct[baseRouter.Path] = pack
		return
	}
	psm.other.PushFront(pack)
}

// find a router with the specific path and return it.
func (psm *Possum) find(path string) (url.Values, HandlerFunc, view.View) {
	defer psm.RUnlock()
	psm.RLock()
	if pack, ok := psm.direct[path]; ok {
		return nil, pack.f, pack.view
	}
	for e := psm.other.Front(); e != nil; e = e.Next() {
		pack := e.Value.(*rvcPack)
		if params, ok := pack.router.Match(path); ok {
			return params, pack.f, pack.view
		}
	}
	return nil, nil, nil
}

// init request handling, and set view and controler function to context.
func (psm *Possum) init(req *http.Request) (HandlerFunc, view.View) {
	params, f, v := psm.find(req.URL.Path)
	if params != nil {
		if (*req).URL.RawQuery == "" {
			(*req).URL.RawQuery = params.Encode()
		} else {
			(*req).URL.RawQuery += "&" + params.Encode()
		}
	}

	if err := (*req).ParseForm(); err != nil {
		panic(Error{http.StatusBadRequest, err.Error()})
	}

	if v == nil || f == nil {
		panic(Error{http.StatusNotFound, "Not Found"})
	}
	return f, v
}

func (psm *Possum) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		r := recover()
		if errPanic, ok := r.(error); ok {
			handleError(w, errPanic)
		}
	}()
	defer psm.preResponseDefer(w, req)
	f, v := psm.init(req)
	psm.preRequest(w, req)
	var data interface{}
	var status int
	if f != nil {
		data, status = f(w, req)
	}
	if link, ok := data.(string); ok {
		if handleRedirect(w, req, status, link) {
			return
		}
	}
	handleRender(v, w, status, data)
}

func (psm *Possum) preRequest(w http.ResponseWriter, req *http.Request) func() {
	return func() {
		if psm.PreRequest != nil {
			psm.PreRequest(w, req)
		}
	}
}

func (psm *Possum) preResponseDefer(w http.ResponseWriter, req *http.Request) func() {
	return func() {
		if psm.PreResponse != nil {
			psm.PreResponse(w, req)
		}
	}
}

func handleRedirect(w http.ResponseWriter, req *http.Request, status int, link string) bool {
	if status != http.StatusMovedPermanently &&
		status != http.StatusFound &&
		status != http.StatusSeeOther &&
		status != http.StatusTemporaryRedirect {
		return false
	}
	http.Redirect(w, req, link, status)
	return true
}

func handleRender(v view.View, w http.ResponseWriter, statusCode int, data interface{}) {
	body, header, err := v.Render(data)
	if err != nil {
		panic(Error{http.StatusInternalServerError, err.Error()})
	}
	if header != nil {
		for key, values := range header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}
	w.WriteHeader(statusCode)
	if _, err = w.Write(body); err != nil {
		panic(Error{http.StatusInternalServerError, err.Error()})
	}
}
