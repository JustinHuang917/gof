package gofcore

import (
	"errors"
	"fmt"
	"reflect"
)

func GetAttributes(m interface{}) (map[string]reflect.Type, error) {
	typ := reflect.TypeOf(m)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	// create an attribute data structure as a map of types keyed by a string.
	attrs := make(map[string]reflect.Type)
	// Only structs are supported so return an empty result if the passed object
	// isn't a struct
	if typ.Kind() != reflect.Struct {
		err_str := fmt.Sprintf("%v type can't have attributes inspected.", typ.Kind())
		err := errors.New(err_str)
		return attrs, err
	}
	// loop through the struct's fields and set the map
	for i := 0; i < typ.NumField(); i++ {
		p := typ.Field(i)
		if !p.Anonymous {
			attrs[p.Name] = p.Type
		}
	}
	return attrs, nil
}

func GetAllMethod(m interface{}) (map[string]reflect.Method, error) {
	typ := reflect.TypeOf(m)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	// create an attribute data structure as a map of types keyed by a string.
	funcs := make(map[string]reflect.Method)
	// Only structs are supported so return an empty result if the passed object
	// isn't a struct
	if typ.Kind() != reflect.Struct {
		err_str := fmt.Sprintf("%v type can't have attributes inspected.", typ.Kind())
		err := errors.New(err_str)
		return funcs, err
	}

	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		funcs[m.Name] = m
	}
	return funcs, nil
}

func CallMethod(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}

func GetFullNameFromType(typ reflect.Type) string {
	if typ == nil {
		return ""
	}
	if typ.Kind() == reflect.Ptr {
		elem := typ.Elem()
		if elem == nil {
			return ""
		} else {
			pkgName := elem.PkgPath()
			typeName := elem.Name()
			return pkgName + "." + typeName
		}
	} else {
		pkgName := typ.PkgPath()
		typeName := typ.Name()
		return pkgName + "." + typeName
	}
	return ""

}
