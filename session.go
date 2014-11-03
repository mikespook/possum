package possum

import "github.com/mikespook/golib/session"

type Session interface {
	Id() string
	Set(key string, value interface{})
	Get(key string) (value interface{})
	Del(key string) (value interface{})
	Init() error
	Clean() error
	Flush() error
}

func StartSession(ctx *Context, f session.FactoryFunc, options session.M) (err error) {
	ctx.Session, err = f(ctx.w, ctx.Request, options)
	return
}
