package possum

import (
	"bytes"
	"encoding/json"

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

type htmlView struct {
	html.Template
}

// This view renders a template into HTML.
func NewHtmlView(filename ...string) (*htmlView, error) {
	t, err := html.ParseFiles(filename...)
	if err != nil {
		return nil, err
	}
	return &htmlView{*t}, nil
}

func (view *htmlView) Render(data interface{}) (output []byte, err error) {
	var buf bytes.Buffer
	if err = view.Execute(&buf, data); err != nil {
		return
	}
	output = buf.Bytes()
	return
}

func (view *htmlView) ContentType() string {
	return "text/html"
}

func (view *htmlView) CharSet() string {
	return "utf-8"
}

type textView struct {
	text.Template
}

// This view renders a template into HTML.
func NewTextView(filename ...string) (*textView, error) {
	t, err := text.ParseFiles(filename...)
	if err != nil {
		return nil, err
	}
	return &textView{*t}, nil
}

func (view *textView) Render(data interface{}) (output []byte, err error) {
	var buf bytes.Buffer
	if err = view.Execute(&buf, data); err != nil {
		return
	}
	output = buf.Bytes()
	return
}

func (view *textView) ContentType() string {
	return "text/plain"
}

func (view *textView) CharSet() string {
	return "utf-8"
}
