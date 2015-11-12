package view

import "net/http"

const (
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html"
	ContentTypePlain  = "text/plain"
	ContentTypeBinary = "application/octet-stream"

	CharSetUTF8 = "utf-8"
)

// View is an interface to render response with a specific format.
type View interface {
	Render(interface{}) ([]byte, http.Header, error)
}
