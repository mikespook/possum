package view

import (
	"bytes"
	"fmt"
	"net/http"
)

// Simple reads and responses data directly.
type simple struct {
	header http.Header
}

func (view simple) Render(data interface{}) (output []byte, err error) {
	if data == nil {
		return nil, nil
	}
	var buf bytes.Buffer
	if _, err = buf.WriteString(fmt.Sprintf("%s", data)); err != nil {
		return
	}
	return buf.Bytes(), nil
}

func (view simple) Header() http.Header {
	return view.header
}

func Simple(contentType, charSet string) simple {
	if contentType == "" {
		contentType = ContentTypePlain
	}
	if charSet == "" {
		charSet = CharSetUTF8
	}
	header := make(http.Header)
	header.Set("Content-Type",
		fmt.Sprintf("%s; charset=%s", contentType, charSet))
	return simple{header}
}
