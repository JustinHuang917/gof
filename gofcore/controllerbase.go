// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofcore

import (
	"bytes"
	"github.com/justinhuang917/gof/gofcore/cfg"
	"net/http"
)

type ControllerBase struct {
}

func (c *ControllerBase) View(v IView, model interface{}, context *HttpContext) (viewResult *ViewResult) {
	viewResult = &ViewResult{Content: new(bytes.Buffer)}
	v.Render(viewResult.Content, model, context.ViewBag, context)
	return
}

func (c *ControllerBase) RedirectToAction(context *HttpContext, actionName string) {
	controllerName := context.ControllerName
	url := cfg.AppConfig.AppPath + "/" + controllerName + "/" + actionName
	http.Redirect(context.ResponseWriter, context.Request, url, 302)
}

func (c *ControllerBase) RedirectToControllerAction(context *HttpContext, controllerName string, actionName string) {
	url := cfg.AppConfig.AppPath + "/" + controllerName + "/" + actionName
	http.Redirect(context.ResponseWriter, context.Request, url, 302)
}
