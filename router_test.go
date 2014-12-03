package possum

import "testing"

func TestResourceRouter(t *testing.T) {
	r := NewResourceRouter("/test/:a/:b/test")
	if !r.Match("/test/a/1/b/2/test") {
		t.Error("Not match!")
		return
	}

	if r.Match("/test/b/1/a/2/test") {
		t.Error("Match!")
		return
	}

	if r.Match("test/a/1/b/2/test") {
		t.Error("Match!")
		return
	}
}
