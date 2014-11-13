package possum

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"path/filepath"
	"sync"
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

// Template is an interface to render template `name` with data,
// and output to wr.
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
	viewWatcher struct {
		*fsnotify.Watcher
		closer chan bool
		count  uint32
	}
)

// InitHtmlTemplates initialzes a series of HTML templates
// in the directory `pattern`.
func InitHtmlTemplates(pattern string) (err error) {
	defer htmlTemp.Unlock()
	htmlTemp.Lock()
	htmlTemp.t, err = html.ParseGlob(pattern)
	return
}

// InitTextTemplates initialzes a series of plain text templates
// in the directory `pattern`.
func InitTextTemplates(pattern string) (err error) {
	defer textTemp.Unlock()
	textTemp.Lock()
	textTemp.t, err = text.ParseGlob(pattern)
	return nil
}

// NewHtmlView retruns a TemplateView witch uses HTML templates internally.
func NewHtmlView(name, charSet string) TemplateView {
	return TemplateView{&htmlTemp, name, "text/html", charSet}
}

// NewTextView retruns a TemplateView witch uses text templates internally.
func NewTextView(name, charSet string) TemplateView {
	return TemplateView{&textTemp, name, "text/plain", charSet}
}

// InitViewWatcher initialzes a watcher to watch templates changes in the `pattern`.
// f would be InitHtmlTemplates or InitTextTemplates.
// If the watcher raises an error internally, the callback function ef will be executed.
// ef can be nil witch represents ignoring all internal errors.
func InitViewWatcher(pattern string, f func(string) error, ef func(error)) (err error) {
	if err = f(pattern); err != nil {
		return
	}
	if viewWatcher.Watcher == nil {
		viewWatcher.Watcher, err = fsnotify.NewWatcher()
		if err != nil {
			return
		}
		viewWatcher.closer = make(chan bool)
	}
	go func() {
		atomic.AddUint32(&viewWatcher.count, 1)
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
			case <-viewWatcher.closer:
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
		if err = viewWatcher.Add(v); err != nil {
			return
		}
	}
	return
}

// CloseViewWatcher closes the wathcer.
func CloseViewWatcher() error {
	for i := uint32(0); i < viewWatcher.count; i++ {
		viewWatcher.closer <- true
	}
	return viewWatcher.Close()
}

// TemplateView represents `html/template` and `text/template` view.
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

// NewStaticFileView returns an new StaticFileView for serving static files.
func NewStaticFileView(filename string, cType string) StaticFileView {
	return StaticFileView{filename, cType}
}

// StaticFileView reads and responses static files.
type StaticFileView struct {
	filename    string
	contentType string
}

func (view StaticFileView) Render(data interface{}) (output []byte, err error) {
	return ioutil.ReadFile(view.filename)
}

func (view StaticFileView) ContentType() string {
	return view.contentType
}

func (view StaticFileView) CharSet() string {
	return ""
}
