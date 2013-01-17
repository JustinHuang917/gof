package gofcore

import (
	"fmt"
	"reflect"
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
	beforeControllerFilterName := getBeforeFilterName("Controller")
	beforeActionFilterName := getBeforeFilterName(actionName)
	beforeContollerFilter := getMethondInController(controller, beforeControllerFilterName)
	beforeActionFilter := getMethondInController(controller, beforeActionFilterName)
	invokeFilter(beforeContollerFilter, context)
	invokeFilter(beforeActionFilter, context)
}

func invokeFilter(filter reflect.Value, context *HttpContext) {
	if filter.IsValid() {
		if filter.Type().NumIn() == 1 {
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
	afterControllerFilterName := getAfterFilterName("Controller")
	afterActionFilterName := getAfterFilterName(actionName)
	afterContollerFilter := getMethondInController(controller, afterControllerFilterName)
	afterActionFilter := getMethondInController(controller, afterActionFilterName)
	invokeFilter(afterActionFilter, context)
	invokeFilter(afterContollerFilter, context)
}

func getBeforeFilterName(identifyName string) string {
	return beforeFilterPrefix + joinChar + identifyName + joinChar + filterEndfix
}

func getAfterFilterName(identifyName string) string {
	return afterFilterPrefix + joinChar + identifyName + joinChar + filterEndfix
}

func getMethondInController(controller interface{}, methondName string) reflect.Value {
	controllerName := reflect.ValueOf(controller).Type().Name()
	key := controllerName + "_" + methondName
	filtersCacheMutex.RLock()
	m, ok := filtersCache[key]
	filtersCacheMutex.RUnlock()
	if !ok {
		filtersCacheMutex.Lock()
		m = reflect.ValueOf(controller).Elem().MethodByName(methondName)
		fmt.Println(m.IsValid())
		filtersCache[key] = m
		filtersCacheMutex.Unlock()
	}
	return m
}
