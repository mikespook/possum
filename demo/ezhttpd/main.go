package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/golib/signal"
	"github.com/mikespook/possum"
	"github.com/mikespook/possum/view"
)

var (
	addr, dir, autoindex string
	errWrongType         = fmt.Errorf("Wrong Type Assertion")
	errAccessDeny        = fmt.Errorf("Access Deny")
)

type staticRouter struct{}

func (r staticRouter) Match(path string) (url.Values, bool) {
	return nil, true
}

type viewData struct {
	contentType string
	body        []byte
}

type staticView struct{}

func (view staticView) Render(data interface{}) (output []byte, err error) {
	if data == nil {
		return nil, errAccessDeny
	}
	switch param := data.(type) {
	case viewData:
		return param.body, nil
	case string:
		return []byte(param), nil
	}
	return nil, errWrongType
}

func (view staticView) Header() http.Header {
	return nil
}

func init() {
	flag.StringVar(&addr, "addr", "127.0.0.1:80", "Served address")
	flag.StringVar(&autoindex, "autoindex", "index.html", "Auto-index file")
	flag.Parse()
}

func main() {
	if dir = flag.Arg(0); dir == "" {
		dir = "."
	}

	mux := possum.NewServerMux()
	mux.HandleFunc(staticRouter{}, newStaticHandle(dir, autoindex), staticView{})
	mux.ErrorHandle = func(err error) {
		log.Error(err)
	}
	mux.PostResponse = func(ctx *possum.Context) error {
		log.Debugf("[%d] %s:%s \"%s\"", ctx.Response.Status,
			ctx.Request.RemoteAddr, ctx.Request.Method,
			ctx.Request.URL.String())
		return nil
	}
	log.Messagef("Addr: %s", addr)
	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Error(err)
			if err := signal.Send(os.Getpid(), os.Interrupt); err != nil {
				panic(err)
			}
		}
	}()
	signal.Bind(os.Interrupt, func() uint {
		log.Message("Exit")
		return signal.BreakExit
	})
	signal.Wait()
}

const tplDir = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Example</title>
<meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
<h1>
<a href="/">ROOT://</a>
{{$current := .current}}
{{range $key, $name := .path}}
<a href="{{$key}}">{{$name}}</a>
{{end}}</h1>
<table>
<thead>
<tr>
	<th></th>
	<th>Name</th>
	<th>Size</th>
	<th>Mode</th>
	<th>Modify Time</th>
</tr>
</thead>
<tbody>
{{range .list}}
<tr>
	<td>{{if .IsDir}}Dir{{end}}</td>
	<td><a href="{{with $current}}{{.}}{{end}}/{{.Name}}">{{.Name}}</a></td>
	<td>{{.Size}}</td>
	<td>{{.Mode}}</td>
	<td>{{.ModTime}}</td>
</tr>
{{end}}
</tbody>
</table>
</body>
</html>`

func newStaticHandle(dir, autoindex string) possum.HandlerFunc {
	return func(ctx *possum.Context) error {
		f := path.Join(dir, ctx.Request.URL.Path)
		fi, err := os.Stat(f)
		if err != nil {
			switch {
			case os.IsNotExist(err):
				ctx.Response.Status = http.StatusNotFound
				ctx.Response.Data = fmt.Sprintf("File Not Found: %s", f)
			case os.IsPermission(err):
				ctx.Response.Status = http.StatusForbidden
				ctx.Response.Data = fmt.Sprintf("Access Forbidden: %s", f)
			default:
				return err
			}
			return nil
		}
		if fi.IsDir() {
			if autoindex != "" {
				ai := path.Join(f, autoindex)
				if _, err := os.Stat(ai); !os.IsNotExist(err) {
					f = ai
					goto F
				}
			}
			fis, err := ioutil.ReadDir(f)
			if err != nil {
				return err
			}
			t, err := template.New("static").Parse(tplDir)
			var buf bytes.Buffer
			current := path.Clean(ctx.Request.URL.Path)
			if current == "/" {
				current = ""
			}
			err = t.Execute(&buf, map[string]interface{}{
				"current": current,
				"path":    splitPath(current),
				"list":    fis,
			})
			if err != nil {
				return err
			}
			ctx.Response.Data = viewData{
				contentType: view.ContentTypeHTML,
				body:        buf.Bytes(),
			}
			return nil
		}
	F:
		body, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		contentType := mime.TypeByExtension(f)
		ctx.Response.Data = viewData{
			contentType: contentType,
			body:        body,
		}
		return nil
	}
}

func splitPath(dir string) (r map[string]string) {
	r = make(map[string]string)
	path := strings.Split(dir, "/")
	key := ""
	for _, name := range path[1:] {
		key = fmt.Sprintf("%s/%s", key, name)
		r[key] = name
	}
	return
}
