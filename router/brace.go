package router

import (
	"net/url"
	"strings"
)

type brace struct {
	matches []string
}

// Brace matches path with REST-full resources URI in brace form.
// e.g., "/foo/v1/bar/v2" will map to "/{foo}/{bar}" form,
// while the value of "foo" and "bar" are "v1" and "v2", respectively.
func Brace(path string) *brace {
	matches := strings.Split(path, "/")
	return &brace{
		matches: matches,
	}
}

func (b *brace) Match(path string) (url.Values, bool) {
	matches := strings.Split(path, "/")
	i := 0
	params := make(url.Values)
	var resKey, resValue string
	for _, v := range b.matches {
		if v != "" && v[0] == '{' && v[len(v)-1] == '}' {
			if matches[i] == v[1:len(v)-1] {
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
