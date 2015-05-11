package view

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

const (
	_body         = "the quick brown fox jumped over the lazy dog"
	_textTemplate = "They said: {{.}}"
	_htmlTemplate = "They said: <p>{{.}}<p>"
)

func _createFile(t *testing.T, content string) (filename string) {
	filename = path.Join(os.TempDir(), "possum.testing")
	if err := ioutil.WriteFile(filename, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return
}

func _deleteFile(filename string) {
	os.Remove(filename)
}
