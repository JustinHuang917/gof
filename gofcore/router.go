package gofcore

import (
	"github.com/JustinHuang917/gof/gofcore/cfg"
	"sort"
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
	return true
}

type argRouteSegement struct {
	ArgName      string
	DefaultValue string
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

func createSegement(seg string) routesegment {
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
		defaultValue := ""
		if l == 1 {
			argName = values[0]
		} else {
			argName = values[0]
			defaultValue = values[1]
		}
		return *&argRouteSegement{ArgName: argName, DefaultValue: defaultValue}
	}
}

type Router struct {
	routeString     string
	routerSegements []routesegment
}

func NewRouter(routeString string) Router {
	segs := splitUrlToPathSegmentString(routeString)
	segements := make([]routesegment, 0, len(segs))
	for _, seg := range segs {
		segements = append(segements, createSegement(seg))
	}
	r := *&Router{routeString: routeString, routerSegements: segements}
	return r
}

func (r *Router) Match(path string) (isMatch bool, args map[string]string) {
	pathParts := splitUrlToPathSegmentString(path)
	l := len(r.routerSegements)
	isMatch = true
	args = make(map[string]string)
	index := 1
	for index, part := range pathParts {
		if index >= l {
			isMatch = false
		}
		segement := r.routerSegements[index]
		if segement != nil {
			isMatch = segement.IsMatch(part)
			if !isMatch {
				break
			} else {
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
		}
	}
	//add default value for not mathed args
	if isMatch && index < l {
		notMatchedSegments := r.routerSegements[index:]
		for _, notMatchedSeg := range notMatchedSegments {
			argSeg, ok := (notMatchedSeg).(argRouteSegement)
			if ok {
				if _, ok := args[argSeg.ArgName]; !ok {
					args[argSeg.ArgName] = argSeg.DefaultValue
				}
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
func (r Routers) Len() int {
	return len(r)
}

func (r Routers) Less(i, j int) bool {
	var iRouteString = r[i].routeString
	var jRouteString = r[j].routeString
	iOrder := cfg.AppConfig.RouteRules[iRouteString]
	jOrder := cfg.AppConfig.RouteRules[jRouteString]
	return iOrder < jOrder
}

func (h Routers) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func RegisterRouter(routeString string) error {
	defer routerSyncMutex.Unlock()
	routerSyncMutex.Lock()
	registeredRouters = append(registeredRouters, NewRouter(routeString))
	sort.Sort(registeredRouters)
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
