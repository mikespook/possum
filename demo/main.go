package main

import (
	"net/http"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/possum"
)

const addr = "127.0.0.1:12345"

func main() {
	mux := possum.NewServerMux()
	mux.HandleFunc("/json", helloworld, possum.JsonView{})

	htmlTemps := possum.NewHtmlTemplates("*.html")
	mux.HandleFunc("/html", helloworld, possum.NewHtmlView(htmlTemps, "base.html"))
	textTemps := possum.NewTextTemplates("*.html")
	mux.HandleFunc("/text", helloworld, possum.NewTextView(textTemps, "base.html"))
	mux.HandleFunc("/project.css", nil, possum.NewFileView("project.css", "text/css"))
	mux.HandleFunc("/img.jpg", nil, possum.NewFileView("img.jpg", "image/jpeg"))
	log.Debug(addr)
	mux.InitPProf("/_pprof")
	http.ListenAndServe(addr, mux)
}

func css(ctx *possum.Context) error {
	return nil
}

func helloworld(ctx *possum.Context) error {
	ctx.Response.Status = http.StatusCreated
	ctx.Response.Data = map[string]interface{}{
		"content": map[string]string{
			"msg":    "hello",
			"target": "world",
		},
	}
	return nil
}
