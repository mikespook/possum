package possum_test

import (
	"fmt"
	"github.com/mikespook/possum"
	"net"
	"net/http"
)

type Foobar struct {
	possum.NoDelete
	possum.NoPatch
	possum.NoPost
	possum.NoPut
}

func (obj *Foobar) Get(w http.ResponseWriter, r *http.Request) (status int, data interface{}) {
	status = http.StatusOK
	data = r.Form
	return
}

func ExampleAddRPC() {
	h := possum.NewHandler()
	h.ErrorHandler = func(err error) {
		fmt.Println(err)
	}
	a := func(w http.ResponseWriter, r *http.Request) (status int, data interface{}) {
		status = http.StatusOK
		data = r.Form
		return
	}
	h.AddRPC("/rpc/test", a)
	http.ListenAndServe(":8080", h)
}

func ExampleNewWrap() {
	h := possum.NewHandler()
	f := func(h possum.HandlerFunc) possum.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) (int, interface{}) {
			if r.Form.Get("secret") != "possum" {
				return http.StatusForbidden, "Wrong secret"
			}
			return h(w, r)
		}
	}
	wrap, err := possum.Wrap(f, &Foobar{})
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
	h.PreHandler = func(r *http.Request) (int, error) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if host != "127.0.0.1" {
			return http.StatusForbidden, fmt.Errorf("Localhost only")
		}
		return http.StatusOK, nil
	}

	h.PostHandler = func(r *http.Request, status int, data interface{}) {
		fmt.Printf("[%d] %s:%s \"%s\"", status, r.RemoteAddr, r.Method, r.URL.String())
	}
	if err := h.AddResource("/rest/test", &Foobar{}); err != nil {
		fmt.Println(err)
		return
	}
	http.ListenAndServe(":8080", h)
}
