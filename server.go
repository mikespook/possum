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
	routers      *Routers
	ErrorHandle  func(error)
	PreRequest   HandlerFunc
	PostResponse HandlerFunc
	NotFound     struct {
		View    View
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
		View    View
		Handler HandlerFunc
	}{SimpleView{}, defaultNotFound}
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
func (mux *ServeMux) HandleFunc(router Router, handler HandlerFunc, view View) {
	router.HandleFunc(handler, view)
	mux.routers.Add(router)
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

func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(w, r)
	router := mux.routers.Find(r.URL.Path)
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				mux.err(e)
			} else {
				mux.err(fmt.Errorf("%s", err))
			}
			return
		}
		var view View
		if router == nil {
			view = mux.NotFound.View
		} else {
			view = router.View()
		}
		if view != nil {
			if err := ctx.flush(view); err != nil {
				mux.err(err)
				return
			}
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
	var handler HandlerFunc
	if router == nil {
		handler = mux.NotFound.Handler
	} else {
		handler = router.Handler()
	}
	if handler != nil {
		if err := handler(ctx); mux.handleError(ctx, err) {
			return
		}
	}
}
