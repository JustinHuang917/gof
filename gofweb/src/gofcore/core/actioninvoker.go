package core

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

func InvokeAction(context *HttpContext) {
	controllerName := context.ControllerName
	controller := GetController(controllerName)
	actionName := getActionName(context, context.ActionName)
	if controller == nil {
		controllerName = context.NoFoundControllerName
		actionName = getActionName(context, context.NoFoundActionName)
		context.RouteName = strings.ToLower("/" + context.NoFoundControllerName + "/" + context.NoFoundActionName)
	}
	controller = GetController(controllerName)
	m := reflect.ValueOf(controller).Elem().MethodByName(actionName)
	if !m.IsValid() {
		controllerName = context.NoFoundControllerName
		actionName = getActionName(context, context.NoFoundActionName)
		context.RouteName = strings.ToLower("/" + context.NoFoundControllerName + "/" + context.NoFoundActionName)
	}
	controller = GetController(controllerName)
	if controller == nil {
		return
	}
	m = reflect.ValueOf(controller).Elem().MethodByName(actionName)
	if !m.IsValid() {
		return
	}
	args := make([]reflect.Value, 1, 1)
	args[0] = reflect.ValueOf(context)
	result := m.Call(args)
	vr := result[0].Interface().(*ViewResult)
	fmt.Fprintf(context.ResponseWriter, string(vr.Content.Bytes()))
}

const (
	getActionPrefix  = "Get"
	postActionPrefix = "Post"
)

func getActionName(context *HttpContext, originalActionName string) string {
	s := originalActionName
	index := 0
	actionName := strings.Map(func(c rune) rune {
		index++
		if index == 1 {
			return unicode.ToUpper(c)
		}

		return unicode.ToLower(c)
	}, s)

	if context.Request.Method == "GET" {
		return getActionPrefix + actionName
	} else if context.Request.Method == "POST" {
		return postActionPrefix + actionName
	}
	return s
}
