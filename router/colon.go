package router

import (
	"net/url"
	"strings"
)

type colon struct {
	matches []string
}

// Colon matches path with REST-full resources URI in perfix colon form.
// e.g., "/foo/v1/bar/v2" will map to "/:foo/:bar" form,
// while the value of "foo" and "bar" are "v1" and "v2", respectively.
func Colon(path string) *colon {
	matches := strings.Split(path, "/")
	return &colon{
		matches: matches,
	}
}

func (c *colon) Match(path string) (url.Values, bool) {
	matches := strings.Split(path, "/")
	i := 0
	params := make(url.Values)
	var resKey, resValue string
	for _, v := range c.matches {
		if v != "" && v[0] == ':' {
			if matches[i] == v[1:] {
				resKey = matches[i]
				i++
				resValue = matches[i]
			} else {
				return nil, false
			}
			params.Add(resKey, resValue)
		} else if matches[i] != v {
			return nil, false
		}
		i++
	}
	return params, true
}
