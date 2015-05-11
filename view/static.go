package view

import "io/ioutil"

// NewStaticFile returns an new StaticFile for serving static files.
func StaticFile(filename string, cType string) staticFile {
	return staticFile{filename, cType}
}

// StaticFile reads and responses static files.
type staticFile struct {
	filename    string
	contentType string
}

func (view staticFile) Render(data interface{}) (output []byte, err error) {
	return ioutil.ReadFile(view.filename)
}

func (view staticFile) ContentType() string {
	return view.contentType
}

func (view staticFile) CharSet() string {
	return ""
}
