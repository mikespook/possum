package possum

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	rpprof "runtime/pprof"
	"strconv"
	"strings"
	"text/template"

	"github.com/mikespook/possum/router"
)

// InitPProf registers pprof handlers to the ServeMux.
// The pprof handlers can be specified a customized prefix.
func (mux *ServerMux) InitPProf(prefix string) {
	if prefix == "" {
		prefix = "/debug/pprof"
	}
	mux.HandleFunc(router.Wildcard(fmt.Sprintf("%s/*", prefix)),
		WrapHTTPHandlerFunc(pprofIndex(prefix)), nil)
	mux.HandleFunc(router.Simple(fmt.Sprintf("%s/cmdline", prefix)),
		WrapHTTPHandlerFunc(http.HandlerFunc(pprof.Cmdline)), nil)
	mux.HandleFunc(router.Simple(fmt.Sprintf("%s/profile", prefix)),
		WrapHTTPHandlerFunc(http.HandlerFunc(pprof.Profile)), nil)
	mux.HandleFunc(router.Simple(fmt.Sprintf("%s/symbol", prefix)),
		WrapHTTPHandlerFunc(http.HandlerFunc(pprof.Symbol)), nil)
}

const pprofTemp = `<html>
<head>
<title>%[1]s/</title>
<style type="text/css">
h1 {border-bottom: 5px solid black;}
</style>
</head>
<body>
<h1>Debug information</h1>
<ul>
	<li><a href="%[1]s/cmdline" target="_blank">Command line</a></li>
	<li><a href="%[1]s/symbol" target="_blank">Symbol</a></li>
	<li><a href="%[1]s/goroutine?debug=2">Full goroutine stack dump</a></li>
</ul>
<h1>Profiles</h1>
<table>
{{range .}}
<tr><td align=right>{{.Count}}<td><a href="%[1]s/{{.Name}}?debug=1">{{.Name}}</a>
{{end}}
<tr><td align=right><td><a href="%[1]s/profile">30-second CPU</a>
</table>
</body>
</html>
`

func pprofIndex(prefix string) http.HandlerFunc {
	var indexTmpl = template.Must(template.New("index").Parse(fmt.Sprintf(pprofTemp, prefix)))
	if prefix[len(prefix)-1] != '/' {
		prefix = fmt.Sprintf("%s/", prefix)
	}
	f := func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, prefix) {
			name := strings.TrimPrefix(r.URL.Path, prefix)
			if name != "" {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				debug, _ := strconv.Atoi(r.FormValue("debug"))
				p := rpprof.Lookup(string(name))
				if p == nil {
					w.WriteHeader(404)
					fmt.Fprintf(w, "Unknown profile: %s\n", name)
					return
				}
				p.WriteTo(w, debug)
				return
			}
		}
		profiles := rpprof.Profiles()
		if err := indexTmpl.Execute(w, profiles); err != nil {
			log.Print(err)
		}
	}
	return f
}
