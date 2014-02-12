package possum

import (
	"fmt"
	"net/http"
	"net/url"
)

// Resources wrap struct
type ResWrap struct{
	f HandlerFunc
	res interface{}
}

// Make a new resources wrap struct for `res` with `f` function
func NewResWrap(f HandlerFunc, res interface{}) (wrap *ResWrap, err error) {
	switch res.(type) {
	case Get, Post, Put, Delete, Patch:
	default:
		return nil, fmt.Errorf("`%T` is not a legal resource", res)
	}
	wrap = &Wrap{f, res}
	return
}

func (wrap *ResWrap) Get(params url.Values) (int, interface{}) {
	if status, err := wrap.f(params); status != http.StatusOK {
		return status, err
	}
	if res, ok := wrap.res.(Get); ok {
		return res.Get(params)
	}
	return (&NoGet{}).Get(params)
}

func (wrap *ResWrap) Post(params url.Values) (int, interface{}) {
	if status, err := wrap.f(params); status != http.StatusOK {
		return status, err
	}
	if res, ok := wrap.res.(Post); ok {
		return res.Post(params)
	}
	return (&NoPost{}).Post(params)
}

func (wrap *ResWrap) Put(params url.Values) (int, interface{}) {
	if status, err := wrap.f(params); status != http.StatusOK {
		return status, err
	}
	if res, ok := wrap.res.(Put); ok {
		return res.Put(params)
	}
	return (&NoPut{}).Put(params)
}

func (wrap *ResWrap) Delete(params url.Values) (int, interface{}) {
	if status, err := wrap.f(params); status != http.StatusOK {
		return status, err
	}
	if res, ok := wrap.res.(Delete); ok {
		return res.Delete(params)
	}
	return (&NoDelete{}).Delete(params)
}

func (wrap *ResWrap) Patch(params url.Values) (int, interface{}) {
	if status, err := wrap.f(params); status != http.StatusOK {
		return status, err
	}
	if res, ok := wrap.res.(Patch); ok {
		return res.Patch(params)
	}
	return (&NoPatch{}).Patch(params)
}
