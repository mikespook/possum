// Possum is a micro web library for Go.
package possum

import (
	"fmt"
	"net/http"
)

// HandlerFunc type is an adapter to allow the use of ordinary
// functions as a HTTP handlers.
type HandlerFunc func(ctx *Context) error

// ServeMux is an HTTP request multiplexer.
type ServeMux struct {
	http.ServeMux
	ErrorHandle  func(error)
	PreRequest   HandlerFunc
	PostResponse HandlerFunc
}

// NewHandler returns a new Handler.
func NewServerMux() (mux *ServeMux) {
	return &ServeMux{*http.NewServeMux(), nil, nil, nil}
}

// Internal error handler
func (mux *ServeMux) err(err error) {
	if mux.ErrorHandle != nil {
		mux.ErrorHandle(err)
	}
}

// HandleFunc specifies a pair of handler and view to handle
// the request witch match pattern.
func (mux *ServeMux) HandleFunc(pattern string, handler HandlerFunc, view View) {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := newContext(w, r)

		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(error); ok {
					mux.err(e)
				} else {
					mux.err(fmt.Errorf("%s", err))
				}
				return
			}
			if err := ctx.flush(view); err != nil {
				mux.err(err)
				return
			}
			if mux.PostResponse != nil {
				if err := mux.PostResponse(ctx); err != nil {
					mux.err(err)
					return
				}
			}
		}()
		if mux.PreRequest != nil {
			if err := mux.PreRequest(ctx); mux.handleError(ctx, err) {
				return
			}
		}
		if handler != nil {
			if err := handler(ctx); mux.handleError(ctx, err) {
				return
			}
		}
	}
	mux.ServeMux.HandleFunc(pattern, f)
}

// handleError tests the context `Error` and assign it to response.
func (mux *ServeMux) handleError(ctx *Context, err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(Error); ok {
		ctx.Response.Status = e.Status
		ctx.Response.Data = e
		return false
	}
	ctx.Response.Status = http.StatusInternalServerError
	ctx.Response.Data = err.Error()
	mux.err(err)
	return true
}
