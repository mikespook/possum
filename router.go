package possum

import (
	"container/list"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/mikespook/possum/router"
	"github.com/mikespook/possum/view"
)

const funcKey contextKey = "func"
const viewKey contextKey = "view"

// Routers contains all routers.
type Routers struct {
	sync.RWMutex
	direct       map[string]*routerPack
	other        *list.List
	PreRequest   http.HandlerFunc
	PostResponse http.HandlerFunc
}

type routerPack struct {
	router router.Router
	view   view.View
	f      HandlerFunc
}

// NewRouters initailizes Routers instance.
func NewRouters() *Routers {
	return &Routers{
		direct: make(map[string]*routerPack),
		other:  list.New(),
	}
}

// Add a router to list
func (routers *Routers) Add(router router.Router, f HandlerFunc, view view.View) {
	defer routers.Unlock()
	routers.Lock()
	pack := &routerPack{router, view, f}
	// direct will full-match the path
	if _, ok := router.(*router.Base); ok {
		routers.direct[sr.Path] = pack
		return
	}
	routers.other.PushFront(pack)
}

// find a router with the specific path and return it.
func (routers *Routers) find(path string) (url.Values, HandlerFunc, view.View) {
	defer routers.RUnlock()
	routers.RLock()
	if pack, ok := routers.direct[path]; ok {
		return nil, pack.f, pack.view
	}
	for e := routers.other.Front(); e != nil; e = e.Next() {
		pack := e.Value.(*routerPack)
		if params, ok := pack.router.Match(path); ok {
			return params, pack.f, pack.view
		}
	}
	return nil, nil, nil
}

func (routers *Routers) init(w http.ResponseWriter, req *http.Request) {
	params, f, view := routers.find(req.URL.Path)
	if params != nil {
		if req.URL.RawQuery == "" {
			req.URL.RawQuery = params.Encode()
		} else {
			req.URL.RawQuery += "&" + params.Encode()
		}
	}

	if err := req.ParseForm(); err != nil {
		panic(Error{http.StatusBadRequest, err.Error()})
	}

	if view == nil || f == nil {
		panic(Errorf{http.StatusNotFound, "Not Found"})
	}
	ctx := req.Context()
	ctx = context.WithValue(ctx, funcKey, f)
	ctx = context.WithValue(ctx, viewKey, view)
	req.WithContext(ctx)
}

func (routers *Routers) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer handleErrorDefer(w)
	defer routers.postRequestDefer(w, req)
	routers.init(w, req)

	handleSession(w, req)
	routers.preRequest(w, req)

	f := getFunc(req.Context())
	if f != nil {
		f(w, req)
	}
	if handleRedirect(w, req) {
		return
	}
	handleRender(w, req)
}

func (routers *Routers) preRequest(w http.ResponseWriter, req *http.Request) {
	return func() {
		if routers.preRequest != nil {
			routers.preRequest(w, req)
		}
	}
}

func (routers *Routers) postRequestDefer(w http.ResponseWriter, req *http.Request) func() {
	return func() {
		if routers.PostRequest != nil {
			routers.PostRequest(w, req)
		}
	}
}
