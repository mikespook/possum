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

// Config struct
type Config struct {
	Addr  string
	PProf string
	Log   configLog
	Test  bool
}

// LoadConfig loads yaml to config instance
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

	psm := possum.New()

	psm.Add(router.Simple("/"), nil, view.Html("index.html", "", ""))
	psm.Add(router.Simple("/json"), helloworld, view.Json("utf-8"))
	psm.Add(router.Wildcard("/json/*/*/*"), helloworld, view.Json("utf-8"))
	psm.Add(router.Simple("/html"), helloworld, view.Html("base.html", "", ""))
	psm.Add(router.Simple("/text"), helloworld, view.Text("base.html", "", ""))
	psm.Add(router.Simple("/project.css"), nil,
		view.StaticFile("statics/project.css", "text/css"))
	tmp, err := view.PreloadFile("statics/img.jpg", "image/jpeg")
	if err != nil {
		log.Error(err)
		return
	}
	psm.Add(router.Simple("/img.jpg"), nil, tmp)
	psm.ErrorHandle = func(err error) {
		log.Error(err)
	}
	psm.PreResponse = func(w http.ResponseWriter, req *http.Request) {
		log.Debugf("[%s] %s \"%s\"", req.Method, req.RemoteAddr, req.URL.String())
	}
	if config.PProf != "" {
		log.Messagef("PProf: http://%s%s", config.Addr, config.PProf)
		possum.InitPProf(psm, config.PProf)
	}
	log.Messagef("Addr: %s", config.Addr)
	go func() {
		if err := http.ListenAndServe(config.Addr, psm); err != nil {
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

func helloworld(w http.ResponseWriter, req *http.Request) (interface{}, int) {
	data := map[string]interface{}{
		"content": map[string]string{
			"msg":    "hello",
			"target": "world",
			"params": req.URL.Query().Encode(),
		},
	}
	w.Header().Set("Test", "Hello world")
	return data, http.StatusCreated
}
