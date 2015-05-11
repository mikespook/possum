package view

import (
	"bytes"
	"fmt"
)

// Simple reads and responses data directly.
type simple struct {
	contentType string
	charSet     string
}

func (view simple) Render(data interface{}) (output []byte, err error) {
	var buf bytes.Buffer
	if _, err = buf.WriteString(fmt.Sprintf("%s", data)); err != nil {
		return
	}
	return buf.Bytes(), nil
}

func (view simple) ContentType() string {
	if view.contentType == "" {

		return ContentTypeHTML
	}
	return view.contentType
}

func (view simple) CharSet() string {
	if view.charSet == "" {
		return CharSetUTF8
	}
	return view.charSet
}

func Simple(contentType, charSet string) simple {
	return simple{contentType, charSet}
}
