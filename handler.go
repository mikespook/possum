// Possum is a micro web-api framework for Go.
package possum

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	PATCH  = "PATCH"
)

// Processing function returning HTTP status code and response object witch will be marshaled to JSON.
type HandlerFunc func(url.Values) (int, interface{})

// Wrapping function
type Wrapper func(w http.ResponseWriter, r *http.Request) (int, interface{})

type (
	// Interface for handling GET request.
	Get interface {
		Get(url.Values) (int, interface{})
	}
	// Interface for handling PUT request.
	Put interface {
		Put(url.Values) (int, interface{})
	}
	// Interface for handling POST request.
	Post interface {
		Post(url.Values) (int, interface{})
	}
	// Interface for handling DELETE request.
	Delete interface {
		Delete(url.Values) (int, interface{})
	}
	// Interface for handling PATCH request.
	Patch interface {
		Patch(url.Values) (int, interface{})
	}
)

// Objects implementing the Get interface response status NotImplemented(501).
type NoGet struct{}
func (NoGet) Get(params url.Values) (int, interface{}) {
	return http.StatusNotImplemented, "GET is not supported"
}

// Objects implementing the Post interface response status NotImplemented(501).
type NoPost   struct{}
func (NoPost) Post(params url.Values) (int, interface{}) {
	return http.StatusNotImplemented, "POST is not supported"
}

// Objects implementing the Put interface response status NotImplemented(501).
type NoPut    struct{}
func (NoPut) Put(params url.Values) (int, interface{}) {
	return http.StatusNotImplemented, "PUT is not supported"
}

// Objects implementing the Delete interface response status NotImplemented(501).
type NoDelete struct{}
func (NoDelete) Delete(params url.Values) (int, interface{}) {
	return http.StatusNotImplemented, "DELETE is not supported"
}

// Objects implementing the Patch interface response status NotImplemented(501).
type NoPatch  struct{}
func (NoPatch) Patch(params url.Values) (int, interface{}) {
	return http.StatusNotImplemented, "PATCH is not supported"
}

// Objects implementing the Handler interface can be registered to serve a particular path or subtree in the HTTP server.
type Hanler struct {
	mux          *http.ServeMux
	ErrorHandler func(error)
	WrapHandler	func(Wrapper) http.HandlerFunc
}

// NewHandler returns a new Handler.
func NewHandler() (h *Hanler) {
	h = &Hanler{
		mux: http.NewServeMux(),
	}
	return
}

// Internal error handler
func (h *Hanler) err(err error) {
	if h.ErrorHandler != nil {
		h.ErrorHandler(err)
	}
}

// Internal wrapping handler
func (h *Hanler) wrap(f Wrapper) http.HandlerFunc {
	if h.WrapHandler == nil {
		return func(w http.ResponseWriter, r *http.Request) {
			f(w, r)
		}
	}
	return h.WrapHandler(f)
}

// AddResource adds a resource to a path. The resource must implement at least one of Get, Post, Put, Delete and Patch interface.
func (h *Hanler) AddResource(pattern string, res interface{}) (err error) {
	switch res.(type) {
	case Get, Post, Put, Delete, Patch:
	default:
		return fmt.Errorf("`%T` is not a legal resource", res)
	}
	h.mux.Handle(pattern, h.wrap(h.rest(res)))
	return
}

// AddRPC adds a Remote Procedure Call to a path.
func (h *Hanler) AddRPC(pattern string, f HandlerFunc) {
	h.mux.Handle(pattern, h.wrap(h.rpc(f)))
	return
}

// Internal wraper for AddResource.
func (h *Hanler) rest(res interface{}) Wrapper {
	return func(w http.ResponseWriter, r *http.Request) (status int, data interface{}) {
		var hf HandlerFunc
		switch r.Method {
		case GET:
			if r, ok := res.(Get); ok {
				hf = r.Get
			} else {
				hf = (&NoGet{}).Get
			}
		case POST:
			if r, ok := res.(Post); ok {
				hf = r.Post
			} else {
				hf = (&NoPost{}).Post
			}
		case PUT:
			if r, ok := res.(Put); ok {
				hf = r.Put
			} else {
				hf = (&NoPut{}).Put
			}
		case DELETE:
			if r, ok := res.(Delete); ok {
				hf = r.Delete
			} else {
				hf = (&NoDelete{}).Delete
			}
		case PATCH:
			if r, ok := res.(Patch); ok {
				hf = r.Patch
			} else {
				hf = (&NoPatch{}).Patch
			}
		}
		f := h.rpc(hf)
		return f(w, r)
	}
}

// Internal wraper for AddRPC.
func (h *Hanler) rpc(f HandlerFunc) Wrapper {
	return func(w http.ResponseWriter, r *http.Request) (status int, data interface{}) {
		status, data = f(r.Form)
		if err := h.writeJson(w, status, data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errstr := err.Error()
			w.Write([]byte(errstr))
			h.err(err)
			return http.StatusInternalServerError, errstr
		}
		return
	}
}

// ServeHTTP calls HandlerFunc
func (h *Hanler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		h.writeErr(w, Errorf(http.StatusBadRequest, "Bad request for `%s`", r.URL.RequestURI()))
		return
	}
	handler, pattern := h.mux.Handler(r)
	if pattern == "" {
		h.writeErr(w, Errorf(http.StatusNotFound, "No handler for `%s`", r.URL.RequestURI()))
		return
	}
	handler.ServeHTTP(w, r)
}

// Internal responsing.
func (h *Hanler) writeJson(w http.ResponseWriter, status int, data interface{}) error {
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	_, err = w.Write(content)
	return err
}

// Internal error responsing.
func (h *Hanler) writeErr(w http.ResponseWriter, apierr apiErr) {
	if err := h.writeJson(w, apierr.status, apierr.message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.err(err)
	}
}
