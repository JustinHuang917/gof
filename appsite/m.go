package main

import (
	"github.com/JustinHuang917/gof/appsite/controller"
	"github.com/JustinHuang917/gof/appsite/view"
	"github.com/JustinHuang917/gof/gofcore"
	"github.com/JustinHuang917/gof/gofcore/cfg"
	"log"
	"net/http"
)

func main() {
	controller.Init()
	view.Init()
	http.HandleFunc(cfg.AppConfig.AppPath, gofcore.Handle)
	err := http.ListenAndServe(cfg.AppConfig.Port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
