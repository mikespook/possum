package main

import (
	"log"
	"net/http"

	"github.com/mikespook/possum"
)

var gitCommit string

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp := possum.NewResponse(r)
		resp.SetData(map[string]string{
			"message": "Hello world",
			"git-commit": gitCommit,
		})
		resp.Write(w)
	})

	log.Println("Starting server on :9000")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatalf("could not start server: %s", err)
	}
}
