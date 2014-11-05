package possum

import (
	"bytes"
	"encoding/json"
	"io"

	html "html/template"
	text "text/template"
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

func NewHtmlTemplates(pattern string) *html.Template {
	return html.Must(html.ParseGlob(pattern))
}

func NewTextTemplates(pattern string) *text.Template {
	return text.Must(text.ParseGlob(pattern))
}

func NewTempView(temp Template, name, cType, charSet string) TempView {
	return TempView{temp, name, cType, charSet}
}

func NewHtmlView(temp Template, name string) TempView {
	return TempView{temp, name, "text/html", "utf-8"}
}

func NewTextView(temp Template, name string) TempView {
	return TempView{temp, name, "text/plain", "utf-8"}
}

type TempView struct {
	Template
	name        string
	contentType string
	charSet     string
}

func (view TempView) Render(data interface{}) (output []byte, err error) {
	var buf bytes.Buffer
	if err = view.ExecuteTemplate(&buf, view.name, data); err != nil {
		return
	}
	output = buf.Bytes()
	return
}

func (view TempView) ContentType() string {
	return view.contentType
}

func (view TempView) CharSet() string {
	return view.charSet
}
