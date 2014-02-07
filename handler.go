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

// Processing function, return HTTP status code and a object
// witch will be marshaled to JSON.
type HandlerFunc func(url.Values) (int, interface{})

type (
	// Get interface for GET request
	Get interface {
		Get(url.Values) (int, interface{})
	}
	// Put interface for PUT request
	Put interface {
		Put(url.Values) (int, interface{})
	}
	// Post interface for POST request
	Post interface {
		Post(url.Values) (int, interface{})
	}
	// Delete interface for DELETE request
	Delete interface {
		Delete(url.Values) (int, interface{})
	}
	// Patch interface for PATCH request
	Patch interface {
		Patch(url.Values) (int, interface{})
	}
)

type (
	NoGet    struct{}
	NoPost   struct{}
	NoPut    struct{}
	NoDelete struct{}
	NoPatch  struct{}
)

func (NoGet) Get(params url.Values) (int, interface{}) {
	return http.StatusNotImplemented, "GET is not supported"
}

func (NoPost) Post(params url.Values) (int, interface{}) {
	return http.StatusNotImplemented, "POST is not supported"
}

func (NoPut) Put(params url.Values) (int, interface{}) {
	return http.StatusNotImplemented, "PUT is not supported"
}

func (NoDelete) Delete(params url.Values) (int, interface{}) {
	return http.StatusNotImplemented, "DELETE is not supported"
}

func (NoPatch) Patch(params url.Values) (int, interface{}) {
	return http.StatusNotImplemented, "PATCH is not supported"
}

type Hanler struct {
	mux          *http.ServeMux
	ErrorHandler func(error)
	PreHandler   func(*http.Request) (int, error)
}

func NewHandler() (h *Hanler) {
	h = &Hanler{
		mux: http.NewServeMux(),
	}
	return
}

func (h *Hanler) err(err error) {
	if h.ErrorHandler != nil {
		h.ErrorHandler(err)
	}
}

func (h *Hanler) AddResource(pattern string, res interface{}) (err error) {
	switch res.(type) {
	case Get, Post, Put, Delete, Patch:
	default:
		return fmt.Errorf("`%T` is not a legal resource", res)
	}
	h.mux.Handle(pattern, h.rest(res))
	return
}

func (h *Hanler) AddRPC(pattern string, f HandlerFunc) (err error) {
	h.mux.Handle(pattern, h.rpc(f))
	return
}

func (h *Hanler) rest(res interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		f(w, r)
	}
}

func (h *Hanler) rpc(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, data := f(r.Form)
		h.writeJson(w, status, data)
	}
}

func (h *Hanler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		h.writeErr(w, newApiErr(http.StatusBadRequest, "Bad request for `%s`", r.URL.RequestURI()))
		return
	}
	if h.PreHandler != nil {
		if status, err := h.PreHandler(r); err != nil {
			h.writeErr(w, newApiErr(status, "%s", err))
			return
		}
	}
	handler, pattern := h.mux.Handler(r)
	if pattern == "" {
		h.writeErr(w, newApiErr(http.StatusNotFound, "No handler for `%s`", r.URL.RequestURI()))
		return
	}
	handler.ServeHTTP(w, r)
}

func (h *Hanler) writeJson(w http.ResponseWriter, status int, data interface{}) {
	content, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(status)
	if _, err := w.Write(content); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (h *Hanler) writeErr(w http.ResponseWriter, apierr apiErr) {
	h.writeJson(w, apierr.status, apierr.message)
}
