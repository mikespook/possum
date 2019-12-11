package possum

import (
	"math/rand"
	"net/http"
)

// Method takes one map as a paramater.
// Keys of this map are HTTP method mapping to HandlerFunc(s).
func Method(m map[string]HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) (interface{}, int) {
		h, ok := m[req.Method]
		if ok {
			return h(w, req)
		}
		return errMethodNotAllowed, http.StatusMethodNotAllowed
	}
}

// Chain combins a slide of HandlerFunc(s) in to one request. TODO
func Chain(h ...HandlerFunc) HandlerFunc {
	f := func(w http.ResponseWriter, req *http.Request) (data interface{}, status int) {
		for _, v := range h {
			data, status = v(w, req)
		}
		return data, status
	}
	return f
}

// Rand picks one HandlerFunc(s) in the slide.
func Rand(h ...HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) (interface{}, int) {
		return h[rand.Intn(len(h))](w, req)
	}
}
