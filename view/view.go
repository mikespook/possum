package view

const (
	ContentTypeJSON  = "application/json"
	ContentTypeHTML  = "text/html"
	ContentTypePlain = "text/plain"

	CharSetUTF8 = "utf-8"
)

// View is an interface to render response with a specific format.
type View interface {
	Render(interface{}) ([]byte, error)
	ContentType() string
	CharSet() string
}
