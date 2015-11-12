package view

import (
	"io/ioutil"
	"mime"
	"net/http"
	"path"
)

// PreloadFile returns the view which can preload static files and serve them.
// The different between `StaticFile` and `PreloadFile` is that `StaticFile`
// load the content of file at every request, while `PreloadFile` load
// the content into memory at the initial stage. Despite that `PreloadFile`
// will be using more memories and could not update the content in time until
// restart the application, it should be fast than `StaticFile` in runtime.
func PreloadFile(filename string, contentType string) (preloadFile, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return preloadFile{}, err
	}
	if contentType == "" {
		contentType = mime.TypeByExtension(path.Ext(filename))
	}
	header := make(http.Header)
	header.Set("Content-Type", contentType)
	return preloadFile{body, header}, nil
}

type preloadFile struct {
	body   []byte
	header http.Header
}

func (view preloadFile) Render(data interface{}) (output []byte, h http.Header, err error) {
	return view.body, view.header, nil
}
