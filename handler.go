// Possum is a micro web-api framework for Go.
package possum

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	PATCH  = "PATCH"

	StatusNone = 0
)

// Processing function returning HTTP status code and response object witch will be marshaled to JSON.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) (int, interface{})

type (
	// Interface for handling GET request.
	Get interface {
		Get(w http.ResponseWriter, r *http.Request) (int, interface{})
	}
	// Interface for handling PUT request.
	Put interface {
		Put(w http.ResponseWriter, r *http.Request) (int, interface{})
	}
	// Interface for handling POST request.
	Post interface {
		Post(w http.ResponseWriter, r *http.Request) (int, interface{})
	}
	// Interface for handling DELETE request.
	Delete interface {
		Delete(w http.ResponseWriter, r *http.Request) (int, interface{})
	}
	// Interface for handling PATCH request.
	Patch interface {
		Patch(w http.ResponseWriter, r *http.Request) (int, interface{})
	}
)

// Objects implementing the Get interface response status NotImplemented(501).
type NoGet struct{}

func (NoGet) Get(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	return http.StatusNotImplemented, "GET is not supported"
}

// Objects implementing the Post interface response status NotImplemented(501).
type NoPost struct{}

func (NoPost) Post(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	return http.StatusNotImplemented, "POST is not supported"
}

// Objects implementing the Put interface response status NotImplemented(501).
type NoPut struct{}

func (NoPut) Put(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	return http.StatusNotImplemented, "PUT is not supported"
}

// Objects implementing the Delete interface response status NotImplemented(501).
type NoDelete struct{}

func (NoDelete) Delete(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	return http.StatusNotImplemented, "DELETE is not supported"
}

// Objects implementing the Patch interface response status NotImplemented(501).
type NoPatch struct{}

func (NoPatch) Patch(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	return http.StatusNotImplemented, "PATCH is not supported"
}

// Objects implementing the Handler interface can be registered to serve a particular path or subtree in the HTTP server.
type Handler struct {
	mux          *http.ServeMux
	ErrorHandler func(error)
	PreHandler   func(r *http.Request) (int, error)
	PostHandler  func(r *http.Request, status int, data interface{})
}

// NewHandler returns a new Handler.
func NewHandler() (h *Handler) {
	h = &Handler{
		mux: http.NewServeMux(),
	}
	return
}

// Internal error handler
func (h *Handler) err(err error) {
	if h.ErrorHandler != nil {
		h.ErrorHandler(err)
	}
}

// Internal wrapping handler
func (h *Handler) wrap(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h.PreHandler != nil {
			status, err := h.PreHandler(r)
			if err != nil {
				h.writeErr(w, Errorf(status, err.Error()))
				return
			}
		}
		status, data := f(w, r)
		if h.PostHandler != nil {
			h.PostHandler(r, status, data)
		}
		switch {
			case status == StatusNone:
			case status >= 300 && status < 400:
				if urlStr, ok := data.(string); ok {
					http.Redirect(w, r, urlStr, status)
				}
			default:
			if err := h.writeJson(w, status, data); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				h.err(err)
			}
		}
	}
}

// AddResource adds a resource to a path. The resource must implement at least one of Get, Post, Put, Delete and Patch interface.
func (h *Handler) AddResource(pattern string, res interface{}) (err error) {
	switch res.(type) {
	case Get, Post, Put, Delete, Patch:
	default:
		return fmt.Errorf("`%T` is not a legal resource", res)
	}
	h.mux.Handle(pattern, h.wrap(h.rest(res)))
	return
}

// AddRPC adds a Remote Procedure Call to a path.
func (h *Handler) AddRPC(pattern string, f HandlerFunc) {
	h.mux.Handle(pattern, h.wrap(f))
}

// Internal wraper for AddResource.
func (h *Handler) rest(res interface{}) HandlerFunc {
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
		return hf(w, r)
	}
}

// ServeHTTP calls HandlerFunc
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
func (h *Handler) writeJson(w http.ResponseWriter, status int, data interface{}) error {
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	_, err = w.Write(content)
	return err
}

// Internal error responsing.
func (h *Handler) writeErr(w http.ResponseWriter, apierr apiErr) {
	if err := h.writeJson(w, apierr.status, apierr.message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.err(err)
	}
}
