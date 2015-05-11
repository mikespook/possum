package view

import j "encoding/json"

// Json renders response data to JSON format.
type json struct {
	charSet string
}

func (view json) Render(data interface{}) (output []byte, err error) {
	return j.Marshal(data)
}

func (view json) ContentType() string {
	return ContentTypeJSON
}

func (view json) CharSet() string {
	if view.charSet == "" {
		return CharSetUTF8
	}
	return view.charSet
}

func Json(charSet string) json {
	return json{charSet}
}
