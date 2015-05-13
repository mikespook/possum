package router

import "testing"

func TestWildcard(t *testing.T) {
	for k, v := range tCases["Wildcard"] {
		r := Wildcard(k)
		testingRouter(t, r, v)
	}
}

func BenchmarkWildcard(b *testing.B) {
	for k, v := range tCases["Wildcard"] {
		r := Wildcard(k)
		benchmarkingRouter(b, r, v)
	}
}
