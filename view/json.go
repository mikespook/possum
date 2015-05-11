package view

import (
	j "encoding/json"
	"fmt"
	"net/http"
)

// Json renders response data to JSON format.
type json struct {
	header http.Header
}

func (view json) Render(data interface{}) (output []byte, err error) {
	return j.Marshal(data)
}

func (view json) Header() http.Header {
	return view.header
}

func Json(charSet string) json {
	if charSet == "" {
		charSet = CharSetUTF8
	}
	header := make(http.Header)
	header.Set("Content-Type",
		fmt.Sprintf("%s; charset=%s", ContentTypeJSON, charSet))
	return json{header}
}
