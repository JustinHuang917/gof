package main

import (
	"fmt"
	"github.com/justinhuang917/gof/appsite/controller"
	"github.com/justinhuang917/gof/appsite/view"
	"github.com/justinhuang917/gof/gofcore"
	"github.com/justinhuang917/gof/gofcore/cfg"
	"log"
	"net/http"
)

func main() {
	// err := cfg.Load("./appsite/cfg.json")
	// if err != nil {
	// 	log.Fatal("Load Config Error: ", err)
	// }
	//gofcore.Init()

	controller.Init()
	view.Init()
	fmt.Println(cfg.AppConfig.AppPath)
	http.HandleFunc(cfg.AppConfig.AppPath, gofcore.Handel)
	err := http.ListenAndServe(cfg.AppConfig.Port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
