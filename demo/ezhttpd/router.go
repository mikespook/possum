package main

import "net/url"

type staticRouter struct{}

func (r staticRouter) Match(path string) (url.Values, bool) {
	return nil, true
}
