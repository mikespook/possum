package view

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestTextTemplates(t *testing.T) {
	f := _createFile(t, _textTemplate)
	defer _deleteFile(f)

	err := InitTextTemplates(os.TempDir() + "/*.testing")
	if err != nil {
		t.Fatal(err)
	}

	v := Text("possum.testing", "", "")
	body, header, err := v.Render(_body)
	if err != nil {
		t.Fatal(err)
	}
	a := strings.Replace(_textTemplate, "{{.}}", _body, -1)
	if string(body) != a {
		t.Fatalf("Rendered template should be %s, got %s.", a, body)
	}
	a = header.Get("Content-Type")
	b := fmt.Sprintf("%s; charset=%s", ContentTypePlain, CharSetUTF8)
	if a != b {
		t.Errorf("Expected Content-Type is %s, got %s.", b, a)
	}
}

func TestHtmlTemplates(t *testing.T) {
	f := _createFile(t, _htmlTemplate)
	defer _deleteFile(f)

	err := InitHtmlTemplates(os.TempDir() + "/*.testing")
	if err != nil {
		t.Fatal(err)
	}

	v := Html("possum.testing", "", "")
	body, header, err := v.Render(_body)
	if err != nil {
		t.Fatal(err)
	}
	a := strings.Replace(_htmlTemplate, "{{.}}", _body, -1)
	if string(body) != a {
		t.Fatalf("Rendered template should be %s, got %s.", _body, a)
	}
	a = header.Get("Content-Type")
	b := fmt.Sprintf("%s; charset=%s", ContentTypeHTML, CharSetUTF8)
	if a != b {
		t.Errorf("Expected Content-Type is %s, got %s.", b, a)
	}
}
