package view

import "testing"

func TestStaticHeader(t *testing.T) {
	f := _createFile(t, _body)
	defer _deleteFile(f)

	v := StaticFile(f, ContentTypeBinary)
	a := v.Header().Get("Content-Type")
	if a != ContentTypeBinary {
		t.Errorf("Expected Content-Type is %s, got %s.", ContentTypeBinary, a)
	}
}

func TestStaticRendering(t *testing.T) {
	f := _createFile(t, _body)
	defer _deleteFile(f)

	sv := StaticFile(f, ContentTypeBinary)
	body, err := sv.Render(nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != _body {
		t.Fatalf("%s should be rendered to %s, got %s.", f, _body, body)
	}
}
