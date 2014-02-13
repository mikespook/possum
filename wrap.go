package possum

import (
	"fmt"
	"net/url"
)

// A wrapper function to HandlerFunc
type WrapFunc func(HandlerFunc) HandlerFunc

// Objects wrapping Resources 
type Wrap struct {
	f   WrapFunc
	res interface{}
}

// NewWrap returns a new Wrap object to add the resource a wrapper function.
func NewWrap(f WrapFunc, res interface{}) (wrap *Wrap, err error) {
	switch res.(type) {
	case Get, Post, Put, Delete, Patch:
	default:
		return nil, fmt.Errorf("`%T` is not a legal resource", res)
	}
	wrap = &Wrap{f, res}
	return
}

// Wrapping Get
func (wrap *Wrap) Get(params url.Values) (int, interface{}) {
	if res, ok := wrap.res.(Get); ok {
		return wrap.f(res.Get)(params)
	}
	return (&NoGet{}).Get(params)
}

// Wrapping Post
func (wrap *Wrap) Post(params url.Values) (int, interface{}) {
	if res, ok := wrap.res.(Post); ok {
		return wrap.f(res.Post)(params)
	}
	return (&NoPost{}).Post(params)
}

// Wrapping Put
func (wrap *Wrap) Put(params url.Values) (int, interface{}) {
	if res, ok := wrap.res.(Put); ok {
		return wrap.f(res.Put)(params)
	}
	return (&NoPut{}).Put(params)
}

// Wrapping Delete
func (wrap *Wrap) Delete(params url.Values) (int, interface{}) {
	if res, ok := wrap.res.(Delete); ok {
		return wrap.f(res.Delete)(params)
	}
	return (&NoDelete{}).Delete(params)
}

// Wrapping Patch
func (wrap *Wrap) Patch(params url.Values) (int, interface{}) {
	if res, ok := wrap.res.(Patch); ok {
		return wrap.f(res.Patch)(params)
	}
	return (&NoPatch{}).Patch(params)
}
