// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofcore

import (
	"bytes"
	"fmt"
	"github.com/justinhuang917/gof/gofcore/cfg"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"sync"
)

type IView interface {
	Render(out *bytes.Buffer, Model interface{}, ViewBag *Bag, httpContext *HttpContext) error
}

type ViewMaps struct {
	mutex *sync.RWMutex
	Views map[string]interface{}
}

// ControllerMaps is a registry for services.

type ControllerMaps struct {
	mutex       *sync.RWMutex
	Controllers map[string]interface{}
}

var viewMaps *ViewMaps

var controllerMaps *ControllerMaps

func RegisterViews(routeName string, view interface{}) {
	fmt.Println("Regeisted View:", routeName)
	viewMaps.mutex.Lock()
	viewMaps.Views[routeName] = view
	viewMaps.mutex.Unlock()
}

func RemoveView(routeName string) {
	viewMaps.mutex.Lock()
	delete(viewMaps.Views, routeName)
	viewMaps.mutex.Unlock()
}

func GetView(routeName string) interface{} {
	defer viewMaps.mutex.RUnlock()
	viewMaps.mutex.RLock()
	return viewMaps.Views[routeName]
}

func RegiesterController(controllerName string, controller interface{}) {
	controllerMaps.mutex.Lock()
	fmt.Println("Regeisted Controller:", controllerName)
	controllerMaps.Controllers[controllerName] = controller
	controllerMaps.mutex.Unlock()
}

func RemoveController(controllerName string) {
	controllerMaps.mutex.Lock()
	delete(controllerMaps.Controllers, controllerName)
	controllerMaps.mutex.Unlock()
}

func GetController(controllerName string) interface{} {
	controllerMaps.mutex.RLock()
	defer controllerMaps.mutex.RUnlock()
	return controllerMaps.Controllers[controllerName]
}

type IHandler interface {
	Handel(context *HttpContext)
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
	DefaultControllerName string
	DefaultActionName     string
	NoFoundControllerName string
	NoFoundActionName     string
	ControllerName        string
	ActionName            string
	Request               *http.Request
	ResponseWriter        http.ResponseWriter
	RouteName             string
	ViewBag               *Bag
	RoutesData            *Bag
	GofSessionId          string
}

func (h *HttpContext) SetSession(key string, value interface{}) {
	SessionMgr.Set(h.GofSessionId, key, value)
}

func (h *HttpContext) GetSession(key string) interface{} {
	return SessionMgr.Get(h.GofSessionId, key)
}

func initHttpContext(w http.ResponseWriter, r *http.Request) *HttpContext {
	url := r.URL
	path := url.Path
	context := new(HttpContext)
	path = strings.TrimLeft(path, cfg.AppConfig.AppPath)
	path = "/" + path
	arr := strings.Split(path, "/")
	l := len(arr)
	if l <= 1 {
		arr = strings.Split("/"+cfg.AppConfig.DefaultPath, "/")
	}
	l = len(arr)
	if l < 3 {
		arr = strings.Split("/"+cfg.AppConfig.NotFoundPath, "/")
	}

	context.ControllerName = strings.ToLower(arr[1])
	context.ActionName = strings.ToLower(arr[2])
	context.DefaultActionName = defaultActionName
	context.DefaultControllerName = defaultControllerName
	context.NoFoundControllerName = noFoundControllerName
	context.NoFoundActionName = noFoundActionName
	context.ResponseWriter = w
	context.Request = r
	context.RouteName = strings.ToLower("/" + context.ControllerName + "/" + context.ActionName)
	context.GofSessionId = ""
	context.ViewBag = NewBag()
	context.RoutesData = NewBag()
	return context
}

type handlerList []IHandler

var handlerListSyncMutex sync.Mutex

func (h handlerList) Len() int {
	return len(h)
}

func (h handlerList) Less(i, j int) bool {
	iName := GetFullNameFromType(reflect.TypeOf(h[i]))
	jName := GetFullNameFromType(reflect.TypeOf(h[j]))
	iOrder := cfg.AppConfig.HandlerSortings[iName]
	jOrder := cfg.AppConfig.HandlerSortings[jName]
	return iOrder < jOrder
	//return h[i].Order() < h[j].Order()
}

func (h handlerList) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

var (
	innerHandlerList      handlerList
	defaultControllerName string
	defaultActionName     string
	noFoundControllerName string
	noFoundActionName     string
)

func RegistHandler(handler IHandler, f func()) {
	handlerListSyncMutex.Lock()
	innerHandlerList = append(innerHandlerList, handler)
	sort.Sort(innerHandlerList)
	if f != nil {
		f()
	}
	handlerListSyncMutex.Unlock()
}

func init() {
	viewMaps = &ViewMaps{new(sync.RWMutex), map[string]interface{}{}}
	controllerMaps = &ControllerMaps{new(sync.RWMutex), map[string]interface{}{}}
	innerHandlerList = make([]IHandler, 0, 10)
	defaultArr := strings.Split(cfg.AppConfig.DefaultPath, "/")
	if len(defaultArr) == 2 {
		defaultControllerName = defaultArr[0]
		defaultActionName = defaultArr[1]
	}
	noFoundArr := strings.Split(cfg.AppConfig.NotFoundPath, "/")
	if len(noFoundArr) == 2 {
		noFoundControllerName = noFoundArr[0]
		noFoundActionName = noFoundArr[1]
	}
	RegistHandler(&DefaultHandler{}, func() {
		fmt.Println("Registered DefaultHandler")
	})
	RegistHandler(&SessionHandler{}, func() {
		fmt.Println("Registered SessionHandler")
		if cfg.AppConfig.EnableSession {
			InitialzieSeesion()
		}
	})

}

func Handel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var context = initHttpContext(w, r)
	context.Request.ParseForm()
	for _, handler := range innerHandlerList {
		if handler != nil {
			handler.Handel(context)
		}
	}
}

type DefaultHandler struct {
}

func (d *DefaultHandler) Handel(context *HttpContext) {
	InvokeAction(context)
}

type SessionHandler struct {
}

func (s *SessionHandler) Handel(context *HttpContext) {
	if cfg.AppConfig.EnableSession {
		sid := cfg.AppConfig.GofSessionId
		ck, err := context.Request.Cookie(sid)
		if err != nil || ck == nil || ck.Value == "" {
			expires := cfg.AppConfig.SessionExpires
			cid, err1 := genUId()
			if err1 == nil {
				c := &http.Cookie{
					Name:  sid,
					Value: cid,
					Path:  "/",
					//Expires:  time.Now().Add(d),
					HttpOnly: true,
					MaxAge:   expires,
				}
				context.Request.AddCookie(c)
				http.SetCookie(context.ResponseWriter, c)
				context.GofSessionId = cid
			} else {
				panic("Generate cookie id error")
			}
		} else {
			context.GofSessionId = ck.Value
		}
	}

}
