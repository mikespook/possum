package possum

import (
	"fmt"
	"net/http"
	"net/http/pprof"
)

func (mux *ServeMux) InitPProf(prefix string) {
	if prefix == "" {
		prefix = "/debug/pprof"
	}
	mux.Handle(fmt.Sprintf("%s/", prefix), http.HandlerFunc(pprof.Index))
	mux.Handle(fmt.Sprintf("%s/cmdline", prefix), http.HandlerFunc(pprof.Cmdline))
	mux.Handle(fmt.Sprintf("%s/profile", prefix), http.HandlerFunc(pprof.Profile))
	mux.Handle(fmt.Sprintf("%s/symbol", prefix), http.HandlerFunc(pprof.Symbol))
}
