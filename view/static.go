package view

import (
	"io/ioutil"
	"net/http"
)

// StaticFile returns the view which can serve static files.
func StaticFile(filename string, contentType string) staticFile {
	if contentType == "" {
		contentType = ContentTypeBinary
	}
	header := make(http.Header)
	header.Set("Content-Type", contentType)
	return staticFile{filename, header}
}

type staticFile struct {
	filename string
	header   http.Header
}

func (view staticFile) Header() http.Header {
	return view.header
}

func (view staticFile) Render(data interface{}) (output []byte, err error) {
	return ioutil.ReadFile(view.filename)
}
