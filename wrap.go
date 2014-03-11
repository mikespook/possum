package possum

import (
	"fmt"
	"net/http"
)

// A wrapper function to HandlerFunc
type wrapFunc func(HandlerFunc) HandlerFunc

// Objects wrapping Resources
type _wrap struct {
	f   wrapFunc
	res interface{}
}

// NewWrap returns a new Wrap object to add the resource a wrapper function.
func Wrap(f wrapFunc, res interface{}) (w *_wrap, err error) {
	switch res.(type) {
	case Get, Post, Put, Delete, Patch:
	default:
		return nil, fmt.Errorf("`%T` is not a legal resource", res)
	}
	w = &_wrap{f, res}
	return
}

// Wrapping Get
func (wrap *_wrap) Get(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	if res, ok := wrap.res.(Get); ok {
		return wrap.f(res.Get)(w, r)
	}
	return (&NoGet{}).Get(w, r)
}

// Wrapping Post
func (wrap *_wrap) Post(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	if res, ok := wrap.res.(Post); ok {
		return wrap.f(res.Post)(w, r)
	}
	return (&NoPost{}).Post(w, r)
}

// Wrapping Put
func (wrap *_wrap) Put(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	if res, ok := wrap.res.(Put); ok {
		return wrap.f(res.Put)(w, r)
	}
	return (&NoPut{}).Put(w, r)
}

// Wrapping Delete
func (wrap *_wrap) Delete(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	if res, ok := wrap.res.(Delete); ok {
		return wrap.f(res.Delete)(w, r)
	}
	return (&NoDelete{}).Delete(w, r)
}

// Wrapping Patch
func (wrap *_wrap) Patch(w http.ResponseWriter, r *http.Request) (int, interface{}) {
	if res, ok := wrap.res.(Patch); ok {
		return wrap.f(res.Patch)(w, r)
	}
	return (&NoPatch{}).Patch(w, r)
}
