package gofcore

import (
	"github.com/JustinHuang917/gof/gofcore/cfg"
	"net/http"
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
	for r, _ := range rules {
		RegisterRouter(r)
	}
}

type RouterHandler struct {
}

func (r *RouterHandler) Handle(context *HttpContext) {
	path := context.Request.URL.Path
	if ok, args := IsMatchRoute(path); ok {
		controerName := args["controller"]
		actionName := args["action"]
		context.ControllerName = controerName
		context.ActionName = actionName
		for k, v := range args {
			context.RoutesData.Add(k, v)
		}
	}
	context.RouteName = strings.ToLower("/" + context.ControllerName + "/" + context.ActionName)
	Debug(context.ControllerName, Runtime)
	Debug(context.ActionName, Runtime)
	Debug(context.RouteName, Runtime)

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
