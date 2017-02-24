// Package possum is a micro web library for Go.
package possum

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/mikespook/possum/router"
	"github.com/mikespook/possum/view"
)

// HandlerFunc type is an adapter to allow the use of ordinary
// functions as a HTTP handlers.
type HandlerFunc func(ctx *Context) error

// ServeMux is an HTTP request multiplexer.
type ServerMux struct {
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

// NewServerMux returns a new Handler.
func NewServerMux() (mux *ServerMux) {
	nf := struct {
		View    view.View
		Handler HandlerFunc
	}{view.Simple(view.ContentTypePlain, view.CharSetUTF8), defaultNotFound}
	return &ServerMux{NewRouters(), nil, nil, nil, nf}
}

// Internal error handler
func (mux *ServerMux) err(err error) {
	if mux.ErrorHandle != nil {
		mux.ErrorHandle(err)
	}
}

// HandleFunc specifies a pair of handler and view to handle
// the request witch matching router.
func (mux *ServerMux) HandleFunc(r router.Router, h HandlerFunc, v view.View) {
	mux.routers.Add(r, h, v)
}

// handleError tests the context `Error` and assign it to response.
func (mux *ServerMux) handleError(ctx *Context, err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(Error); ok {
		ctx.Response.Status = e.Status
		ctx.Response.Data = e
		return true
	}
	if ctx.Response.Status == http.StatusOK {
		ctx.Response.Status = http.StatusInternalServerError
	}
	ctx.Response.Data = err.Error()
	mux.err(err)
	return true
}

func (mux *ServerMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := ctxPool.Get().(*Context)
	ctx.Request = req
	ctx.Response.Status = http.StatusOK
	ctx.Response.ResponseWriter = w

	p, h, v := mux.routers.Find(req.URL.Path)
	if p != nil {
		if ctx.Request.URL.RawQuery == "" {
			ctx.Request.URL.RawQuery = p.Encode()
		} else {
			var buf bytes.Buffer
			buf.WriteString(ctx.Request.URL.RawQuery)
			buf.WriteByte('&')
			buf.WriteString(p.Encode())
			ctx.Request.URL.RawQuery = buf.String()
		}
	}
	if err := ctx.Request.ParseForm(); err != nil {
		mux.err(err)
		return
	}
	if v == nil && h == nil && p == nil {
		h = mux.NotFound.Handler
		v = mux.NotFound.View
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
		if v != nil {
			if err := ctx.flush(v); err != nil {
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
		ctxPool.Put(ctx)
	}()
	if mux.PreRequest != nil {
		if err := mux.PreRequest(ctx); mux.handleError(ctx, err) {
			return
		}
	}
	if h != nil {
		if err := h(ctx); mux.handleError(ctx, err) {
			return
		}
	}
}
