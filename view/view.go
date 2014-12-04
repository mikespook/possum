package view

import (
	"bytes"
	j "encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"sync/atomic"

	html "html/template"
	text "text/template"

	"gopkg.in/fsnotify.v1"
)

// View is an interface to render response with a specific format.
type View interface {
	Render(interface{}) ([]byte, error)
	ContentType() string
	CharSet() string
}

// Json renders response data to JSON format.
type json struct{}

func (view json) Render(data interface{}) (output []byte, err error) {
	return j.Marshal(data)
}

func (view json) ContentType() string {
	return "application/json"
}

func (view json) CharSet() string {
	return "utf-8"
}

func Json() json {
	return json{}
}

// Template is an interface to render template `name` with data,
// and output to wr.
type tmp interface {
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
}

var (
	htmlTemp *html.Template
	textTemp *text.Template
	watcher  struct {
		*fsnotify.Watcher
		closer chan bool
		count  uint32
	}
)

// InitHtmlTemplates initialzes a series of HTML templates
// in the directory `pattern`.
func InitHtmlTemplates(pattern string) (err error) {
	htmlTemp, err = html.ParseGlob(pattern)
	return
}

// InitTextTemplates initialzes a series of plain text templates
// in the directory `pattern`.
func InitTextTemplates(pattern string) (err error) {
	textTemp, err = text.ParseGlob(pattern)
	return nil
}

// NewHtml retruns a TemplateView witch uses HTML templates internally.
func Html(name, charSet string) template {
	return template{htmlTemp, name, "text/html", charSet}
}

// NewText retruns a TemplateView witch uses text templates internally.
func Text(name, charSet string) template {
	return template{textTemp, name, "text/plain", charSet}
}

// InitWatcher initialzes a watcher to watch templates changes in the `pattern`.
// f would be InitHtmlTemplates or InitTextTemplates.
// If the watcher raises an error internally, the callback function ef will be executed.
// ef can be nil witch represents ignoring all internal errors.
func InitWatcher(pattern string, f func(string) error, ef func(error)) (err error) {
	if err = f(pattern); err != nil {
		return
	}
	if watcher.Watcher == nil {
		watcher.Watcher, err = fsnotify.NewWatcher()
		if err != nil {
			return
		}
		watcher.closer = make(chan bool)
	}
	go func() {
		atomic.AddUint32(&watcher.count, 1)
		for {
			select {
			case <-watcher.Events:
				if err := f(pattern); err != nil {
					ef(err)
				}
			case err := <-watcher.Errors:
				if ef != nil {
					ef(err)
				}
			case <-watcher.closer:
				break
			}
		}
	}()

	var matches []string
	matches, err = filepath.Glob(pattern)
	if err != nil {
		return
	}
	for _, v := range matches {
		if err = watcher.Add(v); err != nil {
			return
		}
	}
	return
}

// CloseWatcher closes the wathcer.
func CloseWatcher() error {
	for i := uint32(0); i < watcher.count; i++ {
		watcher.closer <- true
	}
	return watcher.Close()
}

// Template represents `html/template` and `text/template` view.
type template struct {
	tmp
	name        string
	contentType string
	charSet     string
}

func (view template) Render(data interface{}) (output []byte, err error) {
	var buf bytes.Buffer
	if err = view.ExecuteTemplate(&buf, view.name, data); err != nil {
		return
	}
	output = buf.Bytes()
	return
}

func (view template) ContentType() string {
	return view.contentType
}

func (view template) CharSet() string {
	return view.charSet
}

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

// Simple reads and responses data directly.
type simple struct{}

func (view simple) Render(data interface{}) (output []byte, err error) {
	var buf bytes.Buffer
	if _, err = buf.WriteString(fmt.Sprintf("%s", data)); err != nil {
		return
	}
	return buf.Bytes(), nil
}

func (view simple) ContentType() string {
	return "text/html"
}

func (view simple) CharSet() string {
	return "utf-8"
}

func Simple() simple {
	return simple{}
}
