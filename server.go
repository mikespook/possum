// Possum is a micro web library for Go.
package possum

import (
	"fmt"
	"net/http"

	"github.com/mikespook/possum/router"
	"github.com/mikespook/possum/view"
)

// HandlerFunc type is an adapter to allow the use of ordinary
// functions as a HTTP handlers.
type HandlerFunc func(ctx *Context) error

// ServeMux is an HTTP request multiplexer.
type ServeMux struct {
	routers      *Routers
	ErrorHandle  func(error)
	PreRequest   HandlerFunc
	PostResponse HandlerFunc
	NotFound     struct {
		View    view.View
		Handler HandlerFunc
	}
}

var defaultNotFound = func(ctx *Context) error {
	ctx.Response.Status = http.StatusNotFound
	ctx.Response.Data = "Not Found"
	return nil
}

// NewHandler returns a new Handler.
func NewServerMux() (mux *ServeMux) {
	nf := struct {
		View    view.View
		Handler HandlerFunc
	}{view.Simple(), defaultNotFound}
	return &ServeMux{NewRouters(), nil, nil, nil, nf}
}

// Internal error handler
func (mux *ServeMux) err(err error) {
	if mux.ErrorHandle != nil {
		mux.ErrorHandle(err)
	}
}

// HandleFunc specifies a pair of handler and view to handle
// the request witch matching router.
func (mux *ServeMux) HandleFunc(r router.Router, h HandlerFunc, v view.View) {
	mux.routers.Add(r, h, v)
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

func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := newContext(w, req)
	p, h, v := mux.routers.Find(req.URL.Path)
	for k, v1 := range p {
		for _, v2 := range v1 {
			ctx.Request.URL.Query().Add(k, v2)
		}
	}
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				mux.err(e)
			} else {
				mux.err(fmt.Errorf("%s", err))
			}
			return
		}
		if v == nil {
			v = mux.NotFound.View
		}
		if err := ctx.flush(v); err != nil {
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
	if h == nil {
		h = mux.NotFound.Handler
	}
	if err := h(ctx); mux.handleError(ctx, err) {
		return
	}
}
