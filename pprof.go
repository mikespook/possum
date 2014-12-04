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
func (mux *ServeMux) InitPProf(prefix string) {
	if prefix == "" {
		prefix = "/debug/pprof"
	}
	mux.HandleFunc(router.Simple{fmt.Sprintf("%s/", prefix)},
		wrapHttpHandlerFunc(pprofIndex(prefix)), nil)
	mux.HandleFunc(router.Simple{fmt.Sprintf("%s/cmdline", prefix)},
		wrapHttpHandlerFunc(http.HandlerFunc(pprof.Cmdline)), nil)
	mux.HandleFunc(router.Simple{fmt.Sprintf("%s/profile", prefix)},
		wrapHttpHandlerFunc(http.HandlerFunc(pprof.Profile)), nil)
	mux.HandleFunc(router.Simple{fmt.Sprintf("%s/symbol", prefix)},
		wrapHttpHandlerFunc(http.HandlerFunc(pprof.Symbol)), nil)
}

const pprofTemp = `<html>
<head>
<title>%[1]s/</title>
</head>
%[1]s/<br>
<br>
<body>
profiles:<br>
<table>
{{range .}}
<tr><td align=right>{{.Count}}<td><a href="%[1]s/{{.Name}}?debug=1">{{.Name}}</a>
{{end}}
</table>
<br>
<a href="%[1]s/goroutine?debug=2">full goroutine stack dump</a><br>
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

func wrapHttpHandlerFunc(f http.HandlerFunc) HandlerFunc {
	newF := func(ctx *Context) error {
		f(ctx.Response.w, ctx.Request)
		return nil
	}
	return newF
}
