package router

import "testing"

type aTCase map[string]bool

var tCases = map[string]map[string]aTCase{
	"Colon": {
		"/test/:a/:b/test": aTCase{
			"/test/a/1/b/2/test": true,
			"/test/b/1/a/2/test": false,
			"test/a/1/b/2/test":  false,
		}},
	"Brace": {
		"/test/{a}/{b}/test": aTCase{
			"/test/a/1/b/2/test": true,
			"/test/b/1/a/2/test": false,
			"test/a/1/b/2/test":  false,
		}},
	"RegEx": {
		"/test/(.*)/test": aTCase{
			"/test/a/1/b/2/test": true,
			"/foo/b/1/a/2/test":  false,
			"test/a/1/b/2/test":  false,
		}},
	"Wildcard": {
		"/test/*/*/test": aTCase{
			"/test/foo/bar/test": true,
			"/foo/b/1/a/2/test":  false,
			"test/a/1/b/2/test":  false,
		},
	},
}

func testingRouter(t *testing.T, r Router, a aTCase) {
	for k, v := range a {
		if params, ok := r.Match(k); ok != v {
			t.Errorf("%v expected %b, got %b", params, v, ok)
		} else {
			t.Log(params)
		}
	}

}

func benchmarkingRouter(b *testing.B, r Router, a aTCase) {
	for i := 0; i < b.N; i++ {
		for k, v := range a {
			if params, ok := r.Match(k); ok != v {
				b.Errorf("%v expected %b, got %b", params, v, ok)
			}
		}
	}
}

func TestColon(t *testing.T) {
	for k, v := range tCases["Colon"] {
		r := Colon(k)
		testingRouter(t, r, v)
	}
}

func TestBrace(t *testing.T) {
	for k, v := range tCases["Brace"] {
		r := Brace(k)
		testingRouter(t, r, v)
	}
}

func TestRegEx(t *testing.T) {
	for k, v := range tCases["RegEx"] {
		r := RegEx(k)
		testingRouter(t, r, v)
	}
}

func TestWildcard(t *testing.T) {
	for k, v := range tCases["Wildcard"] {
		r := Wildcard(k)
		testingRouter(t, r, v)
	}
}

func BenchmarkColon(b *testing.B) {
	for k, v := range tCases["Colon"] {
		r := Colon(k)
		benchmarkingRouter(b, r, v)
	}
}

func BenchmarkBrace(b *testing.B) {
	for k, v := range tCases["Brace"] {
		r := Brace(k)
		benchmarkingRouter(b, r, v)
	}
}

func BenchmarkRegEx(b *testing.B) {
	for k, v := range tCases["RegEx"] {
		r := RegEx(k)
		benchmarkingRouter(b, r, v)
	}
}

func BenchmarkWildcard(b *testing.B) {
	for k, v := range tCases["Wildcard"] {
		r := Wildcard(k)
		benchmarkingRouter(b, r, v)
	}
}
