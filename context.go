package possum

import "net/http"

type Context struct {
	w       http.ResponseWriter
	Request *http.Request
	Data    interface{}
	Status  int
	Session Session
}

func (ctx *Context) Header() http.Header {
	return ctx.w.Header()
}

func (ctx *Context) flush(view View) error {
	if ctx.Session != nil {
		ctx.Session.Flush()
	}
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
	}
	return ctx
}
