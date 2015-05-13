package router

import "testing"

func TestColon(t *testing.T) {
	for k, v := range tCases["Colon"] {
		r := Colon(k)
		testingRouter(t, r, v)
	}
}

func BenchmarkColon(b *testing.B) {
	for k, v := range tCases["Colon"] {
		r := Colon(k)
		benchmarkingRouter(b, r, v)
	}
}
