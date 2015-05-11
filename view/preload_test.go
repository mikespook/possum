package view

import "testing"

func TestPreloadHeader(t *testing.T) {
	f := _createFile(t, _body)
	defer _deleteFile(f)

	v, err := PreloadFile(f, ContentTypeBinary)
	if err != nil {
		t.Fatal(err)
	}
	a := v.Header().Get("Content-Type")
	if a != ContentTypeBinary {
		t.Errorf("Expected Content-Type is %s, got %s.", ContentTypeBinary, a)
	}
}

func TestPreloadRendering(t *testing.T) {
	f := _createFile(t, _body)
	defer _deleteFile(f)

	pv, err := PreloadFile(f, ContentTypeBinary)
	if err != nil {
		t.Fatal(err)
	}

	body, err := pv.Render(nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != _body {
		t.Fatalf("%s should be rendered to %s, got %s.", f, _body, body)
	}
}
