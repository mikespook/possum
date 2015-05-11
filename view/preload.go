package view

import (
	"io/ioutil"
	"net/http"
)

// PreloadFile returns the view which can preload static files and serve them.
func PreloadFile(filename string, contentType string) (preloadFile, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return preloadFile{}, err
	}
	if contentType == "" {
		contentType = ContentTypeBinary
	}
	header := make(http.Header)
	header.Set("Content-Type", contentType)
	return preloadFile{body, header}, nil
}

type preloadFile struct {
	body   []byte
	header http.Header
}

func (view preloadFile) Header() http.Header {
	return view.header
}

func (view preloadFile) Render(data interface{}) (output []byte, err error) {
	return view.body, nil
}
