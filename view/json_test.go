package view

import (
	"fmt"
	"testing"
)

var jsonTestCases = map[string]interface{}{
	"123":               123,
	"\"foobar\"":        "foobar",
	"true":              true,
	"[0,1,2,3]":         []int{0, 1, 2, 3},
	"{\"Foo\":\"bar\"}": struct{ Foo string }{"bar"},
}

func TestJsonRendering(t *testing.T) {
	jv := Json("")
	for k, v := range jsonTestCases {
		body, header, err := jv.Render(v)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != k {
			t.Fatalf("%v should be rendered to %s, got %s.", v, k, body)
		}
		a := header.Get("Content-Type")
		b := fmt.Sprintf("%s; charset=%s", ContentTypeJSON, CharSetUTF8)
		if a != b {
			t.Errorf("Expected Content-Type is %s, got %s.", b, a)
		}
	}
}
