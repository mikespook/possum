package possum

import (
	"fmt"
	"net/http"

	"github.com/mikespook/possum/session"
)

type Response struct {
	Status int
	Data   interface{}
}

type Context struct {
	w        http.ResponseWriter
	Request  *http.Request
	Response Response
	Session  *session.Session
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
	if charSet != "" {
		charSet = fmt.Sprintf("; charset=%s", charSet)
	}
	ctx.Header().Set("Content-Type", fmt.Sprintf("%s%s", cType, charSet))
	ctx.w.WriteHeader(ctx.Response.Status)
	data, err := view.Render(ctx.Response.Data)
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
		Response: Response{
			Status: http.StatusOK,
		},
	}
	return ctx
}
