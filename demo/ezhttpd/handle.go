package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/mikespook/possum"
	"github.com/mikespook/possum/view"
)

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
		ctx.Response.Data = viewData{
			contentType: mime.TypeByExtension(path.Ext(f)),
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
