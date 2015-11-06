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
	"Simple": {
		"/test/foo/bar/test": aTCase{
			"/test/foo/bar/test": true,
			"/foo/b/1/a/2/test":  false,
			"test/a/1/b/2/test":  false,
		},
	},
}

func testingRouter(t *testing.T, r Router, a aTCase) {
	for k, v := range a {
		if params, ok := r.Match(k); ok != v {
			t.Errorf("%v expected %t, got %t", params, v, ok)
		} else {
			t.Log(params)
		}
	}

}

func benchmarkingRouter(b *testing.B, r Router, a aTCase) {
	for i := 0; i < b.N; i++ {
		for k, v := range a {
			if params, ok := r.Match(k); ok != v {
				b.Errorf("%v expected %t, got %t", params, v, ok)
			}
		}
	}
}

func TestSimple(t *testing.T) {
	for k, v := range tCases["Simple"] {
		r := Simple(k)
		testingRouter(t, r, v)
	}
}

func BenchmarkSimple(b *testing.B) {
	for k, v := range tCases["Simple"] {
		r := Simple(k)
		benchmarkingRouter(b, r, v)
	}
}
