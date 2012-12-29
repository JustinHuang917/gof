package gofcore

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func InvokeAction(context *HttpContext) {
	controllerName := context.ControllerName
	controller := GetController(controllerName)
	actionName := getActionName(context, context.ActionName)
	var result []reflect.Value
	fmt.Println("actionName:", actionName)
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
	if m.Type().NumIn() > 0 {
		args := make([]reflect.Value, 1, 1)
		args[0] = reflect.ValueOf(context)
		if m.Type().NumIn() == 2 {
			kv := convertUrlValuesToMap(context.Request.Form)
			arg := m.Type().In(1)
			sv := reflect.New(arg)
			instance, _ := bindModel(kv, sv, "", -1)
			args = append(args, reflect.ValueOf(instance))
		}
		result = m.Call(args)
	} else {
		result = m.Call(nil)
	}
	vr := result[0].Interface().(*ViewResult)
	fmt.Fprintf(context.ResponseWriter, string(vr.Content.Bytes()))
}

const (
	getActionPrefix  = "Get"
	postActionPrefix = "Post"
)

//Match the action name with action method name
func getActionName(context *HttpContext, originalActionName string) string {
	s := originalActionName
	index := 0
	//First char to upercase eg:action=>Action
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
	fmt.Println(s)
	return s

}

func convertUrlValuesToMap(values url.Values) map[string]string {
	m := make(map[string]string, 0)
	for k, vs := range values {
		if len(vs) > 0 {
			m[k] = vs[0]
		}
	}
	return m
}

func bindModel(kv map[string]string, sv reflect.Value, prefix string, arrayIndex int) (interface{}, bool) {
	svElem := sv.Elem()
	flag := true
	typeOfSV := svElem.Type()
	for i := 0; i < svElem.NumField(); i++ {
		flag = true
		f := svElem.Field(i)
		argName := typeOfSV.Field(i).Name
		var k string
		if prefix == "" {
			k = argName
			//v = kv[argName]
		} else {
			k = prefix + "." + argName
		}
		if arrayIndex >= 0 {
			k = k + "[" + strconv.Itoa(arrayIndex) + "]"
		}
		v, ok := kv[k]
		if !ok {
			flag = false
		} //else {
		kind := typeOfSV.Field(i).Type.Kind()
		//fmt.Println(kind)
		switch kind {
		case reflect.String:
			s := v
			f.SetString(s)
			break
		case reflect.Int:
			n, _ := strconv.ParseInt(v, 10, 64)
			f.SetInt(n)
			break
		case reflect.Array, reflect.Slice:
			keyPrefix := k + "["
			kvs1 := make(map[string]string, 0)
			for k1, v1 := range kv {
				if strings.Index(k1, keyPrefix) == 0 {
					kvs1[k1] = v1
				}
			}
			slice := reflect.MakeSlice(typeOfSV.Field(i).Type, len(kvs1), len(kvs1))
			for k1, v1 := range kvs1 {
				indexStr := strings.TrimLeft(k1, keyPrefix)
				indexStr = strings.TrimRight(indexStr, "]")
				index, _ := strconv.Atoi(indexStr)
				slice.Index(index).Set(reflect.ValueOf(v1))
			}
			f.Set(slice)
			break
		case reflect.Struct:
			_sv := reflect.New(f.Type())
			x, _ := bindModel(kv, _sv, k, -1)
			f.Set(reflect.ValueOf(x))
		}
	}
	sv = reflect.Indirect(sv)
	return sv.Interface(), flag
}