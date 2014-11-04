package possum

import (
	"fmt"
	"net/http"

	"github.com/mikespook/possum/session"
)

type Context struct {
	w       http.ResponseWriter
	Request *http.Request
	Data    interface{}
	Status  int
	Session *session.Session
}

func (ctx *Context) Header() http.Header {
	return ctx.w.Header()
}

func (ctx *Context) flush(view View) error {
	if ctx.Session != nil {
		ctx.Session.Flush()
	}
	cType := view.ContentType()
	if cType == "" {
		cType = "text/plain"
	}
	charSet := view.CharSet()
	if charSet == "" {
		charSet = "utf-8"
	}
	ctx.Header().Set("Content-Type", fmt.Sprintf("%s; charset=%s", cType, charSet))
	ctx.w.WriteHeader(ctx.Status)
	data, err := view.Render(ctx.Data)
	if err != nil {
		return err
	}
	_, err = ctx.w.Write(data)
	return err
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{
		w:       w,
		Request: r,
		Status:  http.StatusOK,
	}
	return ctx
}
