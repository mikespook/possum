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

	htmlView, err := possum.NewHtmlView("base.html", "content.html")
	if err != nil {
		log.Error(err)
		return
	}
	mux.HandleFunc("/html", helloworld, htmlView)

	textView, err := possum.NewTextView("base.html", "content.html")
	if err != nil {
		log.Error(err)
		return
	}
	mux.HandleFunc("/text", helloworld, textView)
	log.Debug(addr)
	http.ListenAndServe(addr, mux)
}

func helloworld(ctx *possum.Context) error {
	ctx.Status = http.StatusCreated
	ctx.Data = map[string]interface{}{
		"content": map[string]string{
			"msg":    "hello",
			"target": "world",
		},
	}
	return nil
}
