package possum_test

import (
	"fmt"
	"github.com/mikespook/possum"
	"net/http"
	"net/url"
)

type Foobar struct {
	possum.NoDelete
	possum.NoPatch
	possum.NoPost
	possum.NoPut
}

func (obj *Foobar) Get(params url.Values) (status int, data interface{}) {
	status = http.StatusOK
	data = params
	return
}

func ExampleAddRPC() {
	h := possum.NewHandler()
	h.ErrorHandler = func(err error) {
		fmt.Println(err)
	}
	a := func(params url.Values) (status int, data interface{}) {
		status = http.StatusOK
		data = params
		return
	}
	h.AddRPC("/rpc/test", a)
	http.ListenAndServe(":8080", h)
}

func ExampleNewWrap() {
	h := possum.NewHandler()
	f := func(h possum.HandlerFunc) possum.HandlerFunc {
		return func(params url.Values) (int, interface{}) {
			if params.Get("secret") != "possum" {
				return http.StatusForbidden, "Wrong secret"
			}
			return h(params)
		}
	}
	wrap, err := possum.Wrap(f , &Foobar{})
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := h.AddResource("/rest/test", wrap); err != nil {
		fmt.Println(err)
		return
	}
	http.ListenAndServe(":8080", h)
}

func ExampleGlobalWrap() {
	h := possum.NewHandler()
	h.WrapHandler = func(h possum.Wrapper) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			status, _ := h(w, r)
			fmt.Printf("[%d] %s:%s \"%s\"", status, r.RemoteAddr, r.Method, r.URL.String())
		}
	}
	if err := h.AddResource("/rest/test", &Foobar{}); err != nil {
		fmt.Println(err)
		return
	}
	http.ListenAndServe(":8080", h)
}
