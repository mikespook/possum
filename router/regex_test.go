package router

import "testing"

func TestRegEx(t *testing.T) {
	for k, v := range tCases["RegEx"] {
		r := RegEx(k)
		testingRouter(t, r, v)
	}
}

func BenchmarkRegEx(b *testing.B) {
	for k, v := range tCases["RegEx"] {
		r := RegEx(k)
		benchmarkingRouter(b, r, v)
	}
}
