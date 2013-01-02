package gofcore

import (
	"fmt"
	"reflect"
	//	"strings"
	"sync"
)

const (
	beforeFilterPrefix = "Before"
	afterFilterPrefix  = "After"
	filterEndfix       = "Filter"
	joinChar           = "_"
)

var (
	filtersCache      map[string]reflect.Value
	filtersCacheMutex *sync.RWMutex
)

func init() {
	filtersCache = make(map[string]reflect.Value, 0)
	filtersCacheMutex = new(sync.RWMutex)
}

func InvokeBeforeFilters(controller interface{}, context *HttpContext, controllerName, actionName string) {
	beforeControllerFilterName := getBeforeFilterName(controllerName)
	beforeActionFilterName := getBeforeFilterName(actionName)
	beforeContollerFilter := getMethondInController(controller, beforeControllerFilterName)
	beforeActionFilter := getMethondInController(controller, beforeActionFilterName)
	invokeFilter(beforeContollerFilter, context)
	invokeFilter(beforeActionFilter, context)

}

func invokeFilter(filter reflect.Value, context *HttpContext) {
	if filter.IsValid() {
		if filter.Type().NumIn() == 0 {
			args := make([]reflect.Value, 1, 1)
			args[0] = reflect.ValueOf(context)
			filter.Call(args)
		} else if filter.Type().NumIn() == 0 {
			filter.Call(nil)
		} else {
			panic("invoke filter error:the number of args more than 1 in filter:" + filter.Type().Name())
		}
	}
}
func InvokeAfterFilters(controller interface{}, context *HttpContext, controllerName, actionName string) {
	afterControllerFilterName := getAfterFilterName(controllerName)
	afterActionFilterName := getAfterFilterName(actionName)
	afterContollerFilter := getMethondInController(controller, afterControllerFilterName)
	afterActionFilter := getMethondInController(controller, afterActionFilterName)
	invokeFilter(afterActionFilter, context)
	invokeFilter(afterContollerFilter, context)
}

func getBeforeFilterName(identifyName string) string {
	identifyName = firstCharToUpper(identifyName)
	return beforeFilterPrefix + joinChar + identifyName + joinChar + filterEndfix
}

func getAfterFilterName(identifyName string) string {
	identifyName = firstCharToUpper(identifyName)
	return afterFilterPrefix + joinChar + identifyName + joinChar + filterEndfix
}

func getMethondInController(controller interface{}, methondName string) reflect.Value {
	filtersCacheMutex.RLock()
	m, ok := filtersCache[methondName]
	filtersCacheMutex.RUnlock()
	if !ok {
		filtersCacheMutex.Lock()
		m = reflect.ValueOf(controller).Elem().MethodByName(methondName)
		filtersCache[methondName] = m
		filtersCacheMutex.Unlock()
	}
	return m
}
