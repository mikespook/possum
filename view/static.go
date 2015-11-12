package view

import (
	"io/ioutil"
	"mime"
	"net/http"
	"path"
)

// StaticFile returns the view which can serve static files.
func StaticFile(filename string, contentType string) staticFile {
	if contentType == "" {
		contentType = mime.TypeByExtension(path.Ext(filename))
	}
	header := make(http.Header)
	header.Set("Content-Type", contentType)
	return staticFile{filename, header}
}

type staticFile struct {
	filename string
	header   http.Header
}

func (view staticFile) Render(data interface{}) (output []byte, h http.Header, err error) {
	output, err = ioutil.ReadFile(view.filename)
	h = view.header
	return
}
