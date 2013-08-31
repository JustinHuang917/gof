// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofcore

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ControllerBase struct {
}

func (c *ControllerBase) RedirectToAction(context *HttpContext, actionName string) {
	url := c.generateUrl(context.ControllerName, actionName, nil)
	http.Redirect(context.ResponseWriter, context.Request, url, 302)
}

func (c *ControllerBase) RedirectToControllerAction(context *HttpContext, controllerName string, actionName string) {
	url := c.generateUrl(controllerName, actionName, nil)
	http.Redirect(context.ResponseWriter, context.Request, url, 302)

}

func (c *ControllerBase) RedirectToActionWithRouteData(context *HttpContext, actionName string, values map[string]string) {
	url := c.generateUrl(context.ControllerName, actionName, values)
	http.Redirect(context.ResponseWriter, context.Request, url, 302)
}
func (c *ControllerBase) RedirectToControllerActionWithRouteData(context *HttpContext, controllerName string, actionName string, values map[string]string) {
	url := c.generateUrl(controllerName, actionName, values)
	http.Redirect(context.ResponseWriter, context.Request, url, 302)
}

func (c *ControllerBase) generateUrl(controllerName string, actionName string, values map[string]string) string {
	return UrlControllerAction(controllerName, actionName, values)
}

func (c *ControllerBase) View(model interface{}, context *HttpContext) (viewResult *ViewResult) {
	v := GetView(context.RouteName).(IView)
	viewResult = &ViewResult{Content: new(bytes.Buffer)}
	err := v.Render(viewResult.Content, model, context.ViewBag, context)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (c *ControllerBase) Json(model interface{}, context *HttpContext) (jsonResult *JsonResult) {
	if bytes, err := json.Marshal(model); err == nil {
		jsonResult = &JsonResult{Content: bytes}
	} else {
		panic(err.Error())
	}
	return jsonResult
}
