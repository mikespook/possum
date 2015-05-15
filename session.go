package possum

import "github.com/mikespook/possum/session"

// StartSession initaillizes a session context.
// This function should be called in a implementation
// of possum.HandleFunc.
func (ctx *Context) StartSession(f session.FactoryFunc) (err error) {
	ctx.Session, err = f(ctx.Response, ctx.Request)
	return
}
