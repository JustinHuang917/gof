package main

import (
	//"bytes"
	"cfg"
	"controller"
	"fmt"
	"gofcore/core"
	"log"
	"net/http"
	//"models"
	"view"
)

func main() {
	//v := core.GetView("/default")
	// if v != nil {
	// 	b := new(bytes.Buffer)
	// 	m := &models.User{"JustinHuang", 0}
	// 	if view, ok := v.(core.IView); ok {
	// 		err := view.Render(b, m)
	// 		if err == nil {
	// 			fmt.Println(string(b.Bytes()))
	// 		} else {
	// 			fmt.Println(err)
	// 		}
	// 	} else {
	// 		fmt.Println("Can't convert to IView")
	// 	}
	// 	fmt.Println(v)
	// } else {
	// 	fmt.Println("No View Found")
	// }
	// dv := &view.V_default{}
	//fmt.Println(dv)
	// c := core.GetController("Home")
	// fmt.Println(c)
	// if c1, ok := c.(*controller.HomeController); ok {
	// 	//b := new(bytes.Buffer)
	// 	//vr := &core.ViewResult{b}
	// 	context := &core.HttpContext{}
	// 	context.ActionName = "Index"
	// 	context.ControllerName = "Home"
	// 	vr := c1.GetIndex(context)
	// 	fmt.Println(string(vr.Content.Bytes()))
	// } else {
	// 	fmt.Println("Not HomeController")
	// }
	// b1 := new(bytes.Buffer)
	// m1 := &models.User{"JustinHuang", 0}
	// err1 := dv.Render(b1, m1)
	// if err1 == nil {
	// 	fmt.Println(string(b1.Bytes()))
	// } else {
	// 	fmt.Println(err1)
	//}
	err := cfg.Load("./cfg.json")
	if err != nil {
		log.Fatal("Load Config Error: ", err)
	}
	view.Init()
	controller.Init()
	fmt.Println(cfg.AppConfig.AppPath)
	http.HandleFunc(cfg.AppConfig.AppPath, core.Handel) //设置访问的路由
	err = http.ListenAndServe(cfg.AppConfig.Port, nil)  //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
