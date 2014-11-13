package possum

import (
	"fmt"
	"net/http"

	"github.com/mikespook/possum/session"
)

// A Response represents an HTTP response status and data to be send to client.
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
	ctx.Response.w.Header().Set("Content-Type", fmt.Sprintf("%s%s", cType, charSet))
	ctx.Response.w.WriteHeader(ctx.Response.Status)
	data, err := view.Render(ctx.Response.Data)
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
