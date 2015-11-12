package main

import (
	"fmt"
	"net/http"

	"github.com/mikespook/possum/view"
)

type viewData struct {
	contentType string
	body        []byte
}

type staticView struct{}

func (v staticView) Render(data interface{}) (output []byte, h http.Header, err error) {
	if data == nil {
		return nil, nil, errAccessDeny
	}
	switch param := data.(type) {
	case viewData:
		header := make(http.Header)
		header.Set("Content-Type", param.contentType)
		return param.body, header, nil
	case string:
		header := make(http.Header)
		header.Set("Content-Type",
			fmt.Sprintf("%s; charset=%s", view.ContentTypePlain, view.CharSetUTF8))
		return []byte(param), header, nil
	}
	return nil, nil, errWrongType
}
