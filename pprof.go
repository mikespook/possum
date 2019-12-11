package possum

import (
	"fmt"
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
func InitPProf(psm *Possum, prefix string) {
	if prefix == "" {
		prefix = "/debug/pprof"
	}
	psm.Add(router.Wildcard(fmt.Sprintf("%s/*", prefix)), pprofIndex(prefix), nil)
	psm.Add(router.Simple(fmt.Sprintf("%s/cmdline", prefix)), convHandlerFunc(pprof.Cmdline), nil)
	psm.Add(router.Simple(fmt.Sprintf("%s/profile", prefix)), convHandlerFunc(pprof.Profile), nil)
	psm.Add(router.Simple(fmt.Sprintf("%s/symbol", prefix)), convHandlerFunc(pprof.Symbol), nil)
}

func convHandlerFunc(f http.HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) (interface{}, int) {
		f(w, req)
		return nil, http.StatusOK
	}
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

func pprofIndex(prefix string) HandlerFunc {
	var indexTmpl = template.Must(template.New("index").Parse(fmt.Sprintf(pprofTemp, prefix)))
	if prefix[len(prefix)-1] != '/' {
		prefix = fmt.Sprintf("%s/", prefix)
	}
	f := func(w http.ResponseWriter, r *http.Request) (interface{}, int) {
		if strings.HasPrefix(r.URL.Path, prefix) {
			name := strings.TrimPrefix(r.URL.Path, prefix)
			if name != "" {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				debug, _ := strconv.Atoi(r.FormValue("debug"))
				p := rpprof.Lookup(string(name))
				if p == nil {
					return fmt.Sprintf("Unknown profile: %s\n", name), http.StatusNotFound
				}
				p.WriteTo(w, debug)
				return nil, http.StatusOK
			}
		}
		profiles := rpprof.Profiles()
		if err := indexTmpl.Execute(w, profiles); err != nil {
			return err, http.StatusInternalServerError
		}
		return nil, http.StatusOK
	}
	return f
}
