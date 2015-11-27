package possum

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/mikespook/possum/session"
	"github.com/mikespook/possum/view"
)

// A Response represents an HTTP response status
// and data to be send to client.
type Response struct {
	Status int
	Data   interface{}
	http.ResponseWriter
}

// A Context contains a Request witch is processed by server handler,
// a Response witch will be send to client
// and a Session witch hold data belonging to a session.
type Context struct {
	Request  *http.Request
	Response Response
	Session  *session.Session
}

// Redirect performs a redirecting to the url, if the code belongs to
// one of StatusMovedPermanently, StatusFound, StatusSeeOther, and
// StatusTemporaryRedirect.
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
			http.Redirect(ctx.Response, ctx.Request, url, ctx.Response.Status)
			return nil
		}
		return NewError(http.StatusInternalServerError,
			fmt.Sprintf("%T is not an URL.", ctx.Response.Data))
	}
	data, header, err := v.Render(ctx.Response.Data)
	if err != nil {
		return err
	}
	if header != nil {
		for hk, hv := range header {
			for _, cv := range hv {
				ctx.Response.Header().Add(hk, cv)
			}
		}
	}
	ctx.Response.WriteHeader(ctx.Response.Status)
	_, err = ctx.Response.Write(data)
	return err
}

var ctxPool *sync.Pool

func init() {
	ctxPool = &sync.Pool{
		New: func() interface{} {
			return &Context{
				Response: Response{},
			}
		},
	}
}
