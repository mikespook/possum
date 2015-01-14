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

// colon matches path with REST-full resources URI in perfix colon form.
type colon struct {
	matches []string
}

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

// brace matches path with REST-full resources URI in brace form.
type brace struct {
	matches []string
}

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
