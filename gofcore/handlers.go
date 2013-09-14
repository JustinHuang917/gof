// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofcore

import (
	"github.com/JustinHuang917/gof/gofcore/cfg"
	"net/http"
	"path"
	"reflect"
	"sort"
	"strings"
	"sync"
)

type handlerList []IHandler

var handlerListSyncMutex *sync.RWMutex

func (h handlerList) Len() int {
	return len(h)
}

func (h handlerList) Less(i, j int) bool {
	iName := GetFullNameFromType(reflect.TypeOf(h[i]))
	jName := GetFullNameFromType(reflect.TypeOf(h[j]))
	iOrder := cfg.AppConfig.HandlerSortings[iName]
	jOrder := cfg.AppConfig.HandlerSortings[jName]
	return iOrder < jOrder
}

func (h handlerList) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

var (
	innerHandlerList handlerList
)

func RegisterHandler(handler IHandler, f func()) {
	handlerListSyncMutex.Lock()
	innerHandlerList = append(innerHandlerList, handler)
	sort.Sort(innerHandlerList)
	if f != nil {
		f()
	}
	handlerListSyncMutex.Unlock()
}

func initHandlers() {
	handlerListSyncMutex = new(sync.RWMutex)
	innerHandlerList = make([]IHandler, 0, 10)
	RegisterHandler(&RouterHandler{}, func() {
		Debug("Registered RouterHandler", StartUp)
	})
	RegisterHandler(&DefaultHandler{}, func() {
		Debug("Registered DefaultHandler", StartUp)
		RegisterRouters()
	})
	RegisterHandler(&SessionHandler{}, func() {
		Debug("Registered SessionHandler", StartUp)
		if cfg.AppConfig.EnableSession {
			InitialzieSeesion()
		}
	})
}

func RegisterRouters() {
	rules := cfg.AppConfig.RouteRules
	for _, r := range rules {
		for k, v := range r {
			RegisterRouter(k, v)
		}
	}
}

type RouterHandler struct {
}

func (r *RouterHandler) Handle(context *HttpContext) {
	p := context.Request.URL.Path
	if ok := r.isStaticContent(context); !ok {
		if ok, args := IsMatchRoute(p); ok {
			controerName := args["controller"]
			actionName := args["action"]
			context.ControllerName = controerName
			context.ActionName = actionName
			for k, v := range args {
				context.RoutesData.Add(k, v)
			}
		}
		//merge query string to route data
		r.mergeRouteDataAndQueryString(context)
		context.RouteName = strings.ToLower("/" + context.ControllerName + "/" + context.ActionName)
		Debug(context.ControllerName, Runtime)
		Debug(context.ActionName, Runtime)
		Debug(context.RouteName, Runtime)
	} else {
		name := path.Join(cfg.AppConfig.RootPath, p)
		http.ServeFile(context.ResponseWriter, context.Request, name)
	}
}
func (r *RouterHandler) mergeRouteDataAndQueryString(context *HttpContext) {
	req := context.Request
	queryValues := req.URL.Query()
	for k, _ := range queryValues {
		//value in querystring will replace value in routermatched
		if k != "controller" && k != "action" {
			context.RoutesData.Add(k, queryValues.Get(k))
		}
	}
}

func (r *RouterHandler) isStaticContent(context *HttpContext) bool {
	req := context.Request
	path := req.URL.Path
	if cfg.AppConfig.StaticDirs != nil {
		for _, dir := range cfg.AppConfig.StaticDirs {
			if hasPrefix := strings.HasPrefix(path, dir); hasPrefix {
				return true
			}
		}
	}
	return false
}

type DefaultHandler struct {
}

func (d *DefaultHandler) Handle(context *HttpContext) {
	InvokeAction(context)
}

type SessionHandler struct {
}

func (s *SessionHandler) Handle(context *HttpContext) {
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
