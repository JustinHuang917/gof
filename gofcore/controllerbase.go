// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofcore

import (
	"bytes"
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
	// url := cfg.AppConfig.AppPath + "/" + controllerName + "/" + actionName
	// if values != nil {
	// 	args := make([]string, 0, len(values))
	// 	for k, v := range values {
	// 		kvStr := fmt.Sprintf("%s=%s", k, v)
	// 		args = append(args, kvStr)
	// 	}
	// 	queryString := strings.Join(args, "&")
	// 	url = fmt.Sprintf("%s?%s", url, queryString)
	// }
	// return url
	return UrlControllerAction(controllerName, actionName, values)
}
