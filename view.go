package possum

import (
	"bytes"
	"encoding/json"
	"text/template"
)

type View interface {
	Render(interface{}) ([]byte, error)
}

type JsonView struct{}

func (view JsonView) Render(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

type htmlView struct {
	template.Template
}

func NewHtmlView(filename ...string) (*htmlView, error) {
	t, err := template.ParseFiles(filename...)
	if err != nil {
		return nil, err
	}
	return &htmlView{*t}, nil
}

func (view *htmlView) Render(data interface{}) (output []byte, err error) {
	var buf bytes.Buffer
	if err := view.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
