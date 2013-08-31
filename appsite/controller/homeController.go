// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"fmt"
	"github.com/JustinHuang917/gof/appsite/models"
	"github.com/JustinHuang917/gof/gofcore"
)

func init() {
	gofcore.RegisterController("home", &HomeController{})
}

type HomeController struct {
	gofcore.ControllerBase
}

func (h HomeController) Before_Controller_Filter(context *gofcore.HttpContext) {
	fmt.Println("filter...")
	v := context.GetSession("username")
	if v == nil && context.ActionName != "login" {
		h.RedirectToAction(context, "login")
		return
	}
}

func (h HomeController) After_Controller_Filter(context *gofcore.HttpContext) {

}

func (h HomeController) Before_GetIndex_Filter(context *gofcore.HttpContext) {

}

func (h HomeController) After_GetIndex_Fitler(context *gofcore.HttpContext) {

}

func (h HomeController) GetIndex(context *gofcore.HttpContext) (viewResult *gofcore.ViewResult) {
	viewResult = h.View(*&models.User{"justinhuang", "", 100, 25}, context)
	return
}

func (h HomeController) PostLogin(context *gofcore.HttpContext, user models.User) (viewResult *gofcore.ViewResult) {
	if user.Name == "justinhuang" && user.Password == "123" {
		context.SetSession("username", "justinhuang")
		h.RedirectToAction(context, "index")
	} else {
		viewResult = h.View(*&models.User{"", "", -1, 0}, context)
	}
	return
}

func (h HomeController) GetLogin(context *gofcore.HttpContext) (viewResult *gofcore.ViewResult) {
	fmt.Println("login...")
	viewResult = h.View(*&models.User{"JustinHuang", "", 100, 25}, context)
	return
}

func (h HomeController) GetNofound(context *gofcore.HttpContext) (viewResult *gofcore.ViewResult) {
	viewResult = h.View(gofcore.NullModel, context)
	return
}
