package gofcore

import (
	//"fmt"
	"github.com/justinhuang917/gof/gofcore/cfg"
	"time"
)

var (
	SessionMgr           sessionManager
	sessionExpires       int
	SessionIsInitialized bool
)

func init() {
	SessionIsInitialized = false
}
func InitialzieSeesion() {
	if !SessionIsInitialized {
		sMode := cfg.AppConfig.SessionMode
		var session ISession
		switch sMode {
		case "InPorc":
			session = *NewInProcSession()
			break
		default:
			session = *NewInProcSession()
		}
		SessionMgr = *&sessionManager{}
		SessionMgr.Session = session
		sessionExpires = cfg.AppConfig.SessionExpires
	}
	SessionIsInitialized = true
}

func setTimeout(timeout time.Duration, name string) {
	t := time.NewTimer(timeout)
	go clearSession(t.C, name)
}

func clearSession(c <-chan time.Time, name string) {
	//SessionMgr.Session.Remove(name)
}

func (s *sessionManager) Get(sessionId, name string) interface{} {
	sname := sessionId + "_" + name
	return s.Session.Get(sname)
}
func (s *sessionManager) Set(sessionId, name string, value interface{}) {
	sname := sessionId + "_" + name
	s.Session.Set(sname, value)
	t := int64(sessionExpires * 1e9)
	td := time.Duration(t)
	setTimeout(td, sname)
}
func (s *sessionManager) Remove(sessionId, name string) {
	sname := sessionId + "_" + name
	s.Session.Remove(sname)
}

type sessionManager struct {
	Session       ISession
	SessionIdName string
}

type ISession interface {
	Get(name string) interface{}
	Set(name string, value interface{})
	Remove(name string)
}
