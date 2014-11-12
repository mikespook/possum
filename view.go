package possum

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"path/filepath"
	"sync"

	html "html/template"
	text "text/template"

	"gopkg.in/fsnotify.v1"
)

// An interface to render response to display.
type View interface {
	Render(interface{}) ([]byte, error)
	ContentType() string
	CharSet() string
}

// JsonView renders response data to JSON format.
type JsonView struct{}

func (view JsonView) Render(data interface{}) (output []byte, err error) {
	return json.Marshal(data)
}

func (view JsonView) ContentType() string {
	return "application/json"
}

func (view JsonView) CharSet() string {
	return "utf-8"
}

type Template interface {
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
}

type tmpTemplate struct {
	sync.Mutex
	t Template
}

func (tmp *tmpTemplate) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	return tmp.t.ExecuteTemplate(wr, name, data)
}

var (
	htmlTemp    tmpTemplate
	textTemp    tmpTemplate
	viewWatcher *fsnotify.Watcher
)

func InitHtmlTemplates(pattern string) (err error) {
	defer htmlTemp.Unlock()
	htmlTemp.Lock()
	htmlTemp.t, err = html.ParseGlob(pattern)
	return
}

func InitTextTemplates(pattern string) (err error) {
	defer textTemp.Unlock()
	textTemp.Lock()
	textTemp.t, err = text.ParseGlob(pattern)
	return nil
}

func NewHtmlView(name, charSet string) TemplateView {
	return TemplateView{&htmlTemp, name, "text/html", charSet}
}

func NewTextView(name, charSet string) TemplateView {
	return TemplateView{&textTemp, name, "text/plain", charSet}
}

func InitViewWatcher(pattern string, f func(string) error, ef func(error)) (err error) {
	if err = f(pattern); err != nil {
		return
	}
	if viewWatcher == nil {
		viewWatcher, err = fsnotify.NewWatcher()
		if err != nil {
			return
		}
		go func() {
			for {
				select {
				case <-viewWatcher.Events:
					if err := f(pattern); err != nil {
						ef(err)
					}
				case err := <-viewWatcher.Errors:
					if ef != nil {
						ef(err)
					}
				}
			}
		}()
	}
	var matches []string
	matches, err = filepath.Glob(pattern)
	if err != nil {
		return
	}
	for _, v := range matches {
		if err = viewWatcher.Add(v); err != nil {
			return
		}
	}
	return
}

func CloseViewWatcher() error {
	return viewWatcher.Close()
}

type TemplateView struct {
	*tmpTemplate
	name        string
	contentType string
	charSet     string
}

func (view TemplateView) Render(data interface{}) (output []byte, err error) {
	var buf bytes.Buffer
	if err = view.ExecuteTemplate(&buf, view.name, data); err != nil {
		return
	}
	output = buf.Bytes()
	return
}

func (view TemplateView) ContentType() string {
	return view.contentType
}

func (view TemplateView) CharSet() string {
	return view.charSet
}

func NewFileView(filename string, cType string) FileView {
	return FileView{filename, cType}
}

type FileView struct {
	filename    string
	contentType string
}

func (view FileView) Render(data interface{}) (output []byte, err error) {
	return ioutil.ReadFile(view.filename)
}

func (view FileView) ContentType() string {
	return view.contentType
}

func (view FileView) CharSet() string {
	return ""
}
