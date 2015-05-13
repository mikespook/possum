package router

import "net/url"

// Router is an interface to match specific path.
type Router interface {
	Match(string) (url.Values, bool)
}

type Base struct {
	Path string
}

// Simple router strictly matches paths.
func Simple(path string) *Base {
	return &Base{path}
}

func (r *Base) Match(path string) (url.Values, bool) {
	return nil, path == r.Path
}
