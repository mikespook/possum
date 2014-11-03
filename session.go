package possum

import "github.com/mikespook/possum/session"

// StartSession initaillizes a session context.
// This function should be called in a implementation
// of possum.HandleFunc.
func StartSession(ctx *Context, f session.FactoryFunc) (err error) {
	ctx.Session, err = f(ctx.w, ctx.Request)
	return
}
