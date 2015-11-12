package view

import (
	"fmt"
	"testing"
)

var simpleTestCases = map[string]interface{}{
	"foobar":          "foobar",
	"<h1>foobar</h1>": "<h1>foobar</h1>",
}

func TestSimpleRendering(t *testing.T) {
	sv := Simple("", "")
	for k, v := range simpleTestCases {
		body, header, err := sv.Render(v)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != k {
			t.Fatalf("%v should be rendered to %s, got %s.", v, k, body)
		}
		a := header.Get("Content-Type")
		b := fmt.Sprintf("%s; charset=%s", ContentTypePlain, CharSetUTF8)
		if a != b {
			t.Errorf("Expected Content-Type is %s, got %s.", b, a)
		}
	}
}
