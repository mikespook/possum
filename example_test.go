package possum_test

import (
	"fmt"
	"github.com/mikespook/possum"
	"net/http"
	"net/url"
)

func a(params url.Values) (status int, data interface{}) {
	status = http.StatusOK
	data = params
	return
}

type b struct {
	possum.NoDelete
	possum.NoPatch
	possum.NoPost
	possum.NoPut
}

func (obj *b) Get(params url.Values) (status int, data interface{}) {
	status = http.StatusOK
	data = params
	return
}

func ExamplePossum() {
	h := possum.NewHandler()
	h.ErrorHandler = func(err error) {
		fmt.Println(err)
	}
	h.PreHandler = func(r *http.Request) (int, error) {
		if r.Form.Get("token") == "" {
			return http.StatusForbidden, fmt.Errorf("Token needed")
		}
		if r.Form.Get("token") == "possum" {
			return http.StatusForbidden, fmt.Errorf("Wrong token")
		}
		return http.StatusOK, nil
	}
	if err := h.AddRPC("/rpc/test", a); err != nil {
		fmt.Println(err)
		return
	}
	if err := h.AddResource("/rest/test", &b{}); err != nil {
		fmt.Println(err)
		return
	}
	http.ListenAndServe(":8080", h)
}
