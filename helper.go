package possum

import (
	"math/rand"
	"net/http"

	"golang.org/x/net/websocket"
)

// Method takes one map as a paramater.
// Keys of this map are HTTP method mapping to HandlerFunc(s).
func Method(m map[string]HandlerFunc) HandlerFunc {
	f := func(ctx *Context) error {
		h, ok := m[ctx.Request.Method]
		if ok {
			return h(ctx)
		}
		ctx.Response.Status = http.StatusMethodNotAllowed
		return nil
	}
	return f
}

// Chain combins a slide of HandlerFunc(s) in to one request.
func Chain(h ...HandlerFunc) HandlerFunc {
	f := func(ctx *Context) error {
		for _, v := range h {
			if err := v(ctx); err != nil {
				return err
			}
		}
		return nil
	}
	return f
}

// Rand picks one HandlerFunc(s) in the slide.
func Rand(h ...HandlerFunc) HandlerFunc {
	f := func(ctx *Context) error {
		if err := h[rand.Intn(len(h))](ctx); err != nil {
			return err
		}
		return nil
	}
	return f
}

// WrapHTTPHandlerFunc wraps http.HandlerFunc in possum.HandlerFunc.
// See pprof.go.
func WrapHTTPHandlerFunc(f http.HandlerFunc) HandlerFunc {
	newF := func(ctx *Context) error {
		f(ctx.Response, ctx.Request)
		return nil
	}
	return newF
}

// WebSocketHandlerFunc convert websocket function to possum.HandlerFunc.
func WebSocketHandlerFunc(f func(ws *websocket.Conn)) HandlerFunc {
	h := websocket.Handler(f)
	return WrapHTTPHandlerFunc(h.ServeHTTP)
}
