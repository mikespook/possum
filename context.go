package possum

import (
	"fmt"
	"net/http"

	"github.com/mikespook/possum/session"
	"github.com/mikespook/possum/view"
)

// A Response represents an HTTP response status
// and data to be send to client.
type Response struct {
	Status int
	Data   interface{}
	w      http.ResponseWriter
}

// Header returns http.Header.
func (response Response) Header() http.Header {
	return response.w.Header()
}

// A Context contains a Request witch is processed by server handler,
// a Response witch will be send to client
// and a Session witch hold data belonging to a session.
type Context struct {
	Request  *http.Request
	Response Response
	Session  *session.Session
}

func (ctx *Context) Redirect(code int, url string) {
	ctx.Response.Status = code
	ctx.Response.Data = url
}

func (ctx *Context) flush(v view.View) error {
	if ctx.Session != nil {
		if err := ctx.Session.Flush(); err != nil {
			return err
		}
	}
	if ctx.Response.Status == http.StatusMovedPermanently ||
		ctx.Response.Status == http.StatusFound ||
		ctx.Response.Status == http.StatusSeeOther ||
		ctx.Response.Status == http.StatusTemporaryRedirect {
		if url, ok := ctx.Response.Data.(string); ok {
			http.Redirect(ctx.Response.w, ctx.Request, url, ctx.Response.Status)
			return nil
		}
		return NewError(http.StatusInternalServerError,
			fmt.Sprintf("%T is not an URL.", ctx.Response.Data))
	}
	cType := v.ContentType()
	if cType == "" {
		cType = "text/plain"
	}
	charSet := v.CharSet()
	if charSet != "" {
		charSet = fmt.Sprintf("; charset=%s", charSet)
	}
	ctx.Response.w.Header().Set("Content-Type",
		fmt.Sprintf("%s%s", cType, charSet))
	ctx.Response.w.WriteHeader(ctx.Response.Status)
	data, err := v.Render(ctx.Response.Data)
	if err != nil {
		return err
	}
	_, err = ctx.Response.w.Write(data)
	return err
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{
		Request: r,
		Response: Response{
			Status: http.StatusOK,
			w:      w,
		},
	}
	return ctx
}
