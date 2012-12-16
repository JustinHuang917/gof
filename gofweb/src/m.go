package main

import (
	"cfg"
	"controller"
	"fmt"
	"gofcore/core"
	"log"
	"net/http"
	"view"
)

func main() {
	err := cfg.Load("./cfg.json")
	if err != nil {
		log.Fatal("Load Config Error: ", err)
	}
	view.Init()
	controller.Init()
	fmt.Println(cfg.AppConfig.AppPath)
	http.HandleFunc(cfg.AppConfig.AppPath, core.Handel)
	err = http.ListenAndServe(cfg.AppConfig.Port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
