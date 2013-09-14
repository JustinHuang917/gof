// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofcore

import (
	"regexp"
	"strings"
	"sync"
)

func isRouteSeprator(str string) bool {
	return str == "/"
}
func isArgRouteSegementString(str string) bool {
	return len(str) > 2 && strings.Index(str, "{") == 0 && strings.LastIndex(str, "}") == len(str)-1
}

func isLiteralSegementString(str string) bool {
	return (!isRouteSeprator(str)) && (!isArgRouteSegementString(str))
}

type routesegment interface {
	IsMatch(str string) bool
}

type sepratorRouteSegment struct {
}
type literalRouteSegment struct {
	str string
}

func (s sepratorRouteSegment) IsMatch(str string) bool {
	return isRouteSeprator(str)
}

func (l literalRouteSegment) IsMatch(str string) bool {
	if l.str == "*" {
		return true
	}
	return strings.EqualFold(strings.ToLower(str), strings.ToLower(l.str))
}

func (a argRouteSegement) IsMatch(str string) bool {
	if a.regex != nil {
		return a.regex.MatchString(str)
	}
	return true
}

type argRouteSegement struct {
	ArgName         string
	RegexString     string
	DefaultValue    string
	hasDefaultValue bool
	regex           *regexp.Regexp
}

func newArgRouteSegement(argName, regex, defaultValue string, hasDefaultValue bool) argRouteSegement {
	a := &argRouteSegement{ArgName: argName,
		RegexString:     regex,
		DefaultValue:    defaultValue,
		hasDefaultValue: hasDefaultValue}
	if regex != "" {
		a.regex = regexp.MustCompile(regex)
	}
	return *a
}
func formatSegemenetString(str string) string {
	s := strings.Replace(str, "{{", "{", 0)
	s = strings.Replace(s, "}}", "}", 0)
	return s
}

func splitUrlToPathSegmentString(url string) []string {
	segs := make([]string, 0, 3)
	if url != "" {
		index := 0
		l := len(url)
		sep := "/"
		for i := 0; i < l; i = index + 1 {
			index = indexOfString(url, sep, i)
			if index == -1 && i < l {
				str := url[i:]
				if len(str) > 0 {
					segs = append(segs, str)
				}
				return segs
			}
			seg := url[i:index]
			if len(seg) > 0 {
				segs = append(segs, seg)
			}
			segs = append(segs, sep)
		}
	}
	return segs
}

func createSegement(seg string, defaultValues map[string]string) routesegment {
	if isRouteSeprator(seg) {
		return *&sepratorRouteSegment{}
	} else if isLiteralSegementString(seg) {
		return *&literalRouteSegment{str: seg}
	} else {
		value := strings.TrimLeft(seg, "{")
		value = strings.TrimRight(value, "}")
		values := strings.Split(value, ":")
		l := len(values)
		argName := ""
		regexString := ""
		if l == 1 {
			argName = values[0]
		} else {
			argName = values[0]
			regexString = values[1]
		}
		defaultValue := ""
		hasDefaultValue := false
		if defaultValues != nil {
			if v, ok := defaultValues[argName]; ok {
				defaultValue = v
				hasDefaultValue = ok
			}
		}
		return newArgRouteSegement(argName, regexString, defaultValue, hasDefaultValue)
	}
}

type Router struct {
	routeString     string
	routerSegements []routesegment
	defaultValues   map[string]string
}

func NewRouter(routeString string, defaultValues map[string]string) Router {
	segs := splitUrlToPathSegmentString(routeString)
	segements := make([]routesegment, 0, len(segs))
	for _, seg := range segs {
		segements = append(segements, createSegement(seg, defaultValues))
	}
	r := *&Router{routeString: routeString,
		routerSegements: segements,
		defaultValues:   defaultValues}
	return r
}

func (r *Router) Match(path string) (isMatch bool, args map[string]string) {
	pathParts := splitUrlToPathSegmentString(path)
	l := len(pathParts)
	isMatch = true
	args = make(map[string]string)
	for index, segement := range r.routerSegements {
		if index < l { //when  parts of path less than router parts
			part := pathParts[index]
			isMatch = segement.IsMatch(part)
			if !isMatch {
				break
			} else { //if is matched and it's argRouteSegement,need add part as value to return
				argSeg, ok := (segement).(argRouteSegement)
				if ok {
					//the matched argname can't be duplicate
					if _, ok := args[argSeg.ArgName]; ok {
						isMatch = false
						break
					} else {
						args[argSeg.ArgName] = part
					}
				}
			}
		} else {
			//pass when is seprator
			if _, ok := (segement).(sepratorRouteSegment); ok {
				isMatch = true
				continue
			}
			//when it's argRouteSegement and has default value, pass, and add default value to return
			argSeg, ok := (segement).(argRouteSegement)
			if ok && argSeg.hasDefaultValue {
				//the matched argname can't be duplicate
				if _, ok := args[argSeg.ArgName]; ok {
					isMatch = false
					break
				} else {
					args[argSeg.ArgName] = argSeg.DefaultValue
				}
			} else {
				isMatch = false
				break
			}
		}
	}

	//add not matched default values to return
	if r.defaultValues != nil && isMatch {
		for k, v := range r.defaultValues {
			if _, ok := args[k]; !ok {
				args[k] = v
			}
		}
	}
	return
}

type Routers []Router

var routerSyncMutex *sync.RWMutex
var registeredRouters Routers

func initRouters() {
	registeredRouters = make([]Router, 0, 5)
	routerSyncMutex = new(sync.RWMutex)
}

func RegisterRouter(routeString string, defaultValues map[string]string) error {
	defer routerSyncMutex.Unlock()
	routerSyncMutex.Lock()
	registeredRouters = append(registeredRouters, NewRouter(routeString, defaultValues))
	// sort.Sort(registeredRouters)
	return nil
}

func IsMatchRoute(path string) (isMatch bool, args map[string]string) {
	defer routerSyncMutex.RUnlock()
	routerSyncMutex.RLock()
	for _, r := range registeredRouters {
		if isMatch, args = r.Match(path); isMatch {
			Debug("match rule"+r.routeString, Runtime)
			break
		}
	}
	return
}
