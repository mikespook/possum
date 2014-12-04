package router

import (
	"net/url"
	"regexp"
	"strings"
)

// Router is an interface to match specific path.
type Router interface {
	Match(string) (url.Values, bool)
}

type Base struct {
	Path string
}

func Simple(path string) Base {
	return Base{path}
}

func (r Base) Match(path string) (url.Values, bool) {
	return nil, path == r.Path
}

// wildcard matches path with wildcard
type wildcard struct {
	matches []string
}

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
	for k, v := range r.matches {
		if v != "*" && matches[k] != v {
			return nil, false
		}
	}
	return nil, true
}

// regex matches path with regex
type regex struct {
	r *regexp.Regexp
}

func RegEx(path string) *regex {
	r, err := regexp.Compile(path)
	if err != nil {
		panic(err)
	}
	return &regex{
		r: r,
	}
}

func (r *regex) Match(path string) (url.Values, bool) {
	return nil, r.r.MatchString(path)
}

// resource matches path with REST-full resources form.
type resource struct {
	matches []string
}

func Resource(path string) *resource {
	matches := strings.Split(path, "/")
	return &resource{
		matches: matches,
	}
}

func (r *resource) Match(path string) (url.Values, bool) {
	matches := strings.Split(path, "/")
	i := 0
	params := make(url.Values)
	var resKey, resValue string
	for _, v := range r.matches {
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
