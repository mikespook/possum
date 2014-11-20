package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/golib/signal"
	"github.com/mikespook/possum"
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
	flag.StringVar(&configFile, "config", "", "Path to the configuration file")
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
		if err := possum.InitViewWatcher("*.html", possum.InitHtmlTemplates, nil); err != nil {
			log.Error(err)
			return
		}
		if err := possum.InitViewWatcher("*.html", possum.InitTextTemplates, nil); err != nil {
			log.Error(err)
			return
		}
	} else {
		if err := possum.InitHtmlTemplates("*.html"); err != nil {
			log.Error(err)
			return
		}

		if err := possum.InitTextTemplates("*.html"); err != nil {
			log.Error(err)
			return
		}
	}

	mux := possum.NewServerMux()
	mux.HandleFunc("/json", helloworld, possum.JsonView{})
	mux.HandleFunc("/html", helloworld, possum.NewHtmlView("base.html", "utf-8"))
	mux.HandleFunc("/text", helloworld, possum.NewTextView("base.html", "utf-8"))
	mux.HandleFunc("/project.css", nil, possum.NewStaticFileView("project.css", "text/css"))
	mux.HandleFunc("/img.jpg", nil, possum.NewStaticFileView("img.jpg", "image/jpeg"))

	if config.PProf != "" {
		log.Messagef("PProf: http://%s%s", config.Addr, config.PProf)
		mux.InitPProf(config.PProf)
	}
	log.Messagef("Addr: %s", config.Addr)
	go http.ListenAndServe(config.Addr, mux)
	sh := signal.NewHandler()
	sh.Bind(os.Interrupt, func() bool {
		log.Message("Exit")
		return true
	})
	sh.Loop()
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
	ctx.Response.Header().Set("Test", "Hello world")
	return nil
}
