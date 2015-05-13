package router

import (
	"net/url"
	"strings"
)

type wildcard struct {
	matches []string
}

// Wildcard matches paths with wildcard form.
// E.g., "/foo/v1/bar/v2" will match the form "/foo/*/bar/*",
// and "/foo/v1/v2/bar" will match the form "/foo/*/*/bar", but
// will not match "/foo/*/bar".
func Wildcard(path string) *wildcard {
	matches := strings.Split(path, "/")
	return &wildcard{
		matches: matches,
	}
}

func (r *wildcard) Match(path string) (url.Values, bool) {
	matches := strings.Split(path, "/")
	if len(matches) != len(r.matches) {
		return nil, false
	}
	p := url.Values{}
	for k, v := range r.matches {
		if v != "*" && matches[k] != v {
			return nil, false
		}
		if v == "*" && matches[k] != "" {
			p.Add(matches[k], matches[k])
		}
	}
	return p, true
}
