package possum

import (
	"fmt"
	"net/url"
)

// A wrapper function to HandlerFunc
type wrapFunc func(HandlerFunc) HandlerFunc

// Objects wrapping Resources
type wrap struct {
	f   wrapFunc
	res interface{}
}

// NewWrap returns a new Wrap object to add the resource a wrapper function.
func Wrap(f wrapFunc, res interface{}) (w *wrap, err error) {
	switch res.(type) {
	case Get, Post, Put, Delete, Patch:
	default:
		return nil, fmt.Errorf("`%T` is not a legal resource", res)
	}
	w = &wrap{f, res}
	return
}

// Wrapping Get
func (w *wrap) Get(params url.Values) (int, interface{}) {
	if res, ok := w.res.(Get); ok {
		return w.f(res.Get)(params)
	}
	return (&NoGet{}).Get(params)
}

// Wrapping Post
func (w *wrap) Post(params url.Values) (int, interface{}) {
	if res, ok := w.res.(Post); ok {
		return w.f(res.Post)(params)
	}
	return (&NoPost{}).Post(params)
}

// Wrapping Put
func (w *wrap) Put(params url.Values) (int, interface{}) {
	if res, ok := w.res.(Put); ok {
		return w.f(res.Put)(params)
	}
	return (&NoPut{}).Put(params)
}

// Wrapping Delete
func (w *wrap) Delete(params url.Values) (int, interface{}) {
	if res, ok := w.res.(Delete); ok {
		return w.f(res.Delete)(params)
	}
	return (&NoDelete{}).Delete(params)
}

// Wrapping Patch
func (w *wrap) Patch(params url.Values) (int, interface{}) {
	if res, ok := w.res.(Patch); ok {
		return w.f(res.Patch)(params)
	}
	return (&NoPatch{}).Patch(params)
}
