package router

import "testing"

func TestBrace(t *testing.T) {
	for k, v := range tCases["Brace"] {
		r := Brace(k)
		testingRouter(t, r, v)
	}
}

func BenchmarkBrace(b *testing.B) {
	for k, v := range tCases["Brace"] {
		r := Brace(k)
		benchmarkingRouter(b, r, v)
	}
}
