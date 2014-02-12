package possum

import (
	"fmt"
	"net/url"
)

type WrapFunc func(HandlerFunc) HandlerFunc

// Resources wrap struct
type Wrap struct {
	f   WrapFunc
	res interface{}
}

// Make a new resources wrap struct for resource `res` with function `f`
func NewWrap(f WrapFunc, res interface{}) (wrap *Wrap, err error) {
	switch res.(type) {
	case Get, Post, Put, Delete, Patch:
	default:
		return nil, fmt.Errorf("`%T` is not a legal resource", res)
	}
	wrap = &Wrap{f, res}
	return
}

func (wrap *Wrap) Get(params url.Values) (int, interface{}) {
	if res, ok := wrap.res.(Get); ok {
		return wrap.f(res.Get)(params)
	}
	return (&NoGet{}).Get(params)
}

func (wrap *Wrap) Post(params url.Values) (int, interface{}) {
	if res, ok := wrap.res.(Post); ok {
		return wrap.f(res.Post)(params)
	}
	return (&NoPost{}).Post(params)
}

func (wrap *Wrap) Put(params url.Values) (int, interface{}) {
	if res, ok := wrap.res.(Put); ok {
		return wrap.f(res.Put)(params)
	}
	return (&NoPut{}).Put(params)
}

func (wrap *Wrap) Delete(params url.Values) (int, interface{}) {
	if res, ok := wrap.res.(Delete); ok {
		return wrap.f(res.Delete)(params)
	}
	return (&NoDelete{}).Delete(params)
}

func (wrap *Wrap) Patch(params url.Values) (int, interface{}) {
	if res, ok := wrap.res.(Patch); ok {
		return wrap.f(res.Patch)(params)
	}
	return (&NoPatch{}).Patch(params)
}
