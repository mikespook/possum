package router

import "testing"

func TestColon(t *testing.T) {
	r := Colon("/test/:a/:b/test")
	if params, ok := r.Match("/test/a/1/b/2/test"); !ok {
		t.Error("Not match!", params)
		return
	} else {
		t.Log(params)
	}

	if params, ok := r.Match("/test/b/1/a/2/test"); ok {
		t.Error("Match!", params)
		return
	} else {
		t.Log(params)
	}

	if params, ok := r.Match("test/a/1/b/2/test"); ok {
		t.Error("Match!", params)
		return
	} else {
		t.Log(params)
	}
}

func TestBrace(t *testing.T) {
	r := Brace("/test/{a}/{b}/test")
	if params, ok := r.Match("/test/a/1/b/2/test"); !ok {
		t.Error("Not match!", params)
		return
	} else {
		t.Log(params)
	}

	if params, ok := r.Match("/test/b/1/a/2/test"); ok {
		t.Error("Match!", params)
		return
	} else {
		t.Log(params)
	}

	if params, ok := r.Match("test/a/1/b/2/test"); ok {
		t.Error("Match!", params)
		return
	} else {
		t.Log(params)
	}
}

func TestRegEx(t *testing.T) {
	r := RegEx("/test/(.*)/test")
	if params, ok := r.Match("/test/a/1/b/2/test"); !ok {
		t.Error("Not match!", params)
		return
	}

	if params, ok := r.Match("/test1/b/1/a/2/test1"); ok {
		t.Error("Match!", params)
		return
	}

	if params, ok := r.Match("test/a/1/b/2/test"); ok {
		t.Error("Match!", params)
		return
	}
}

func TestWildcard(t *testing.T) {
	r := Wildcard("/test/*/*/test")
	if params, ok := r.Match("/test/a/b/test"); !ok {
		t.Error("Not match!", params)
		return
	}

	if params, ok := r.Match("/test/b/1/a/2/test"); ok {
		t.Error("Match!", params)
		return
	}

	if params, ok := r.Match("/a/1/b/2/test"); ok {
		t.Error("Match!", params)
		return
	}
}
