// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofcore

import (
	"bytes"
	"net/http"
	"strings"
	"sync"
)

var (
	defaultControllerName string
	defaultActionName     string
	noFoundControllerName string
	noFoundActionName     string
	viewMaps              *ViewMaps
	controllerMaps        *ControllerMaps
)

type IView interface {
	Render(out *bytes.Buffer, Model interface{}, ViewBag *Bag, httpContext *HttpContext) error
}

type ViewMaps struct {
	mutex *sync.RWMutex
	Views map[string]interface{}
}

type ControllerMaps struct {
	mutex       *sync.RWMutex
	Controllers map[string]interface{}
}

func RegisterViews(routeName string, view interface{}) {
	viewMaps.mutex.Lock()
	routeName = strings.ToLower(routeName)
	Debug("Regeisted View:"+routeName, StartUp)
	viewMaps.Views[routeName] = view
	viewMaps.mutex.Unlock()
}

func RemoveView(routeName string) {
	viewMaps.mutex.Lock()
	routeName = strings.ToLower(routeName)
	routeName = strings.ToLower(routeName)
	delete(viewMaps.Views, routeName)
	viewMaps.mutex.Unlock()
}

func GetView(routeName string) interface{} {
	defer viewMaps.mutex.RUnlock()
	viewMaps.mutex.RLock()
	routeName = strings.ToLower(routeName)
	Debug("routename:"+routeName, Runtime)
	return viewMaps.Views[routeName]
}

func RegisterController(controllerName string, controller interface{}) {
	controllerMaps.mutex.Lock()
	controllerName = strings.ToLower(controllerName)
	Debug("Regeisted Controller:"+controllerName, StartUp)
	controllerMaps.Controllers[controllerName] = controller
	controllerMaps.mutex.Unlock()
}

func RemoveController(controllerName string) {
	controllerMaps.mutex.Lock()
	controllerName = strings.ToLower(controllerName)
	delete(controllerMaps.Controllers, controllerName)
	controllerMaps.mutex.Unlock()
}

func GetController(controllerName string) interface{} {
	controllerMaps.mutex.RLock()
	defer controllerMaps.mutex.RUnlock()
	controllerName = strings.ToLower(controllerName)
	return controllerMaps.Controllers[controllerName]
}

type IHandler interface {
	Handle(context *HttpContext)
}
type Bag struct {
	Bags  map[string]interface{}
	mutex *sync.RWMutex
}

func NewBag() *Bag {
	bag := new(Bag)
	bag.Bags = make(map[string]interface{})
	bag.mutex = new(sync.RWMutex)
	return bag
}

func (v *Bag) Add(key string, value interface{}) {
	v.mutex.Lock()
	v.Bags[key] = value
	v.mutex.Unlock()
}

func (v *Bag) Remove(key string) {
	v.mutex.Lock()
	delete(v.Bags, key)
	v.mutex.Unlock()
}

func (v *Bag) Get(key string) interface{} {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return v.Bags[key]
}

type HttpContext struct {
	ControllerName string
	ActionName     string
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	RouteName      string
	ViewBag        *Bag
	RoutesData     *Bag
	GofSessionId   string
}

func (h *HttpContext) SetSession(key string, value interface{}) {
	SessionMgr.Set(h.GofSessionId, key, value)
}

func (h *HttpContext) GetSession(key string) interface{} {
	return SessionMgr.Get(h.GofSessionId, key)
}

func initHttpContext(w http.ResponseWriter, r *http.Request) *HttpContext {
	context := new(HttpContext)
	context.ResponseWriter = w
	context.Request = r
	context.GofSessionId = ""
	context.ViewBag = NewBag()
	context.RoutesData = NewBag()
	return context
}

func initApplication() {
	viewMaps = &ViewMaps{new(sync.RWMutex), map[string]interface{}{}}
	controllerMaps = &ControllerMaps{new(sync.RWMutex), map[string]interface{}{}}
}

func Handle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var context = initHttpContext(w, r)
	context.Request.ParseForm()
	for _, handler := range innerHandlerList {
		if handler != nil {
			handler.Handle(context)
		}
	}
}
