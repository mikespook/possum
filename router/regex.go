package router

import (
	"net/url"
	"regexp"
)

type regex struct {
	r *regexp.Regexp
}

// RegEx matches path using regular patterns.
// TODO Dumpping sub-patterns into url.Values.
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
