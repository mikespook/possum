package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/golib/signal"
	"github.com/mikespook/possum"
)

var (
	addr, dir, autoindex string
	errWrongType         = fmt.Errorf("Wrong Type Assertion")
	errAccessDeny        = fmt.Errorf("Access Deny")
)

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
