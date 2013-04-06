// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	//"bytes"
	//"fmt"
	"github.com/justinhuang917/gof/appsite/models"
	"github.com/justinhuang917/gof/gofcore"
)

func init() {
	gofcore.RegiesterController("home", &HomeController{})
}

type HomeController struct {
	gofcore.ControllerBase
}

func (h HomeController) Before_Controller_Filter(context *gofcore.HttpContext) {
	cid := context.GofSessionId
	v := gofcore.SessionMgr.Get(cid, "username")
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
	v1 := gofcore.GetView(context.RouteName)
	viewResult = gofcore.View(v1.(gofcore.IView), &models.User{"justinhuang", "", 100, 25}, context)
	return
}

func (h HomeController) PostLogin(context *gofcore.HttpContext, user models.User) (viewResult *gofcore.ViewResult) {
	if user.Name == "justinhuang" && user.Password == "123" {
		cid := context.GofSessionId
		gofcore.SessionMgr.Set(cid, "username", "justinhuang")
		h.RedirectToAction(context, "index")
	} else {
		v := gofcore.GetView(context.RouteName)
		viewResult = gofcore.View(v.(gofcore.IView), &models.User{"", "", -1, 0}, context)
	}
	return
}

func (h HomeController) GetLogin(context *gofcore.HttpContext) (viewResult *gofcore.ViewResult) {
	v := gofcore.GetView(context.RouteName)
	viewResult = gofcore.View(v.(gofcore.IView), &models.User{"JustinHuang", "", 100, 25}, context)
	return
}

func (h HomeController) GetNofound(context *gofcore.HttpContext) (viewResult *gofcore.ViewResult) {
	v := gofcore.GetView(context.RouteName)
	viewResult = gofcore.View(v.(gofcore.IView), gofcore.NullModel, context)
	return
}
