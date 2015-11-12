package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/golib/signal"
	"github.com/mikespook/possum"
	"github.com/mikespook/possum/router"
	"github.com/mikespook/possum/view"
	"gopkg.in/yaml.v2"
)

type configLog struct {
	File, Level string
}

type Config struct {
	Addr  string
	PProf string
	Log   configLog
	Test  bool
}

func LoadConfig(filename string) (config *Config, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &config)
	return
}

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "Path to the configuration file")
	flag.Parse()
}

func main() {
	if configFile == "" {
		flag.Usage()
		return
	}
	config, err := LoadConfig(configFile)
	if err != nil {
		log.Error(err)
		flag.Usage()
		return
	}
	if err := log.Init(config.Log.File, log.StrToLevel(config.Log.Level), log.DefaultCallDepth); err != nil {
		log.Error(err)
	}

	if config.Test {
		if err := view.InitWatcher("./templates/*.html", view.InitHtmlTemplates, nil); err != nil {
			log.Error(err)
			return
		}
		if err := view.InitWatcher("./templates/*.html", view.InitTextTemplates, nil); err != nil {
			log.Error(err)
			return
		}
	} else {
		if err := view.InitHtmlTemplates("./templates/*.html"); err != nil {
			log.Error(err)
			return
		}

		if err := view.InitTextTemplates("./templates/*.html"); err != nil {
			log.Error(err)
			return
		}
	}

	mux := possum.NewServerMux()

	mux.HandleFunc(router.Simple("/"), nil, view.Html("index.html", "", ""))
	mux.HandleFunc(router.Simple("/json"), helloworld, view.Json("utf-8"))
	mux.HandleFunc(router.Wildcard("/json/*/*/*"), helloworld, view.Json("utf-8"))
	mux.HandleFunc(router.Simple("/html"), helloworld, view.Html("base.html", "", ""))
	mux.HandleFunc(router.Simple("/text"), helloworld, view.Text("base.html", "", ""))
	mux.HandleFunc(router.Simple("/project.css"), nil, view.StaticFile("statics/project.css", "text/css"))
	tmp, err := view.PreloadFile("statics/img.jpg", "image/jpeg")
	if err != nil {
		log.Error(err)
		return
	}
	mux.HandleFunc(router.Simple("/img.jpg"), nil, tmp)
	mux.ErrorHandle = func(err error) {
		log.Error(err)
	}
	mux.PostResponse = func(ctx *possum.Context) error {
		log.Debugf("[%d] %s:%s \"%s\"", ctx.Response.Status,
			ctx.Request.RemoteAddr, ctx.Request.Method,
			ctx.Request.URL.String())
		return nil
	}
	if config.PProf != "" {
		log.Messagef("PProf: http://%s%s", config.Addr, config.PProf)
		mux.InitPProf(config.PProf)
	}
	log.Messagef("Addr: %s", config.Addr)
	go func() {
		if err := http.ListenAndServe(config.Addr, mux); err != nil {
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

func css(ctx *possum.Context) error {
	return nil
}

func helloworld(ctx *possum.Context) error {
	ctx.Response.Status = http.StatusCreated
	ctx.Response.Data = map[string]interface{}{
		"content": map[string]string{
			"msg":    "hello",
			"target": "world",
			"params": ctx.Request.URL.Query().Encode(),
		},
	}
	ctx.Response.Header().Set("Test", "Hello world")
	return nil
}
