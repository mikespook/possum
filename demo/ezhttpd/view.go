package main

import "net/http"

type viewData struct {
	contentType string
	body        []byte
}

type staticView struct{}

func (view staticView) Render(data interface{}) (output []byte, err error) {
	if data == nil {
		return nil, errAccessDeny
	}
	switch param := data.(type) {
	case viewData:
		return param.body, nil
	case string:
		return []byte(param), nil
	}
	return nil, errWrongType
}

func (view staticView) Header() http.Header {
	return nil
}
