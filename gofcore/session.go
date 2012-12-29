package gofcore

// import (
// 	"sync"
// 	"github.com/justinhuang917/gof/gofcore/cfg"
// )

// type SessionManager *SessionManager

// func func init() {
// 	SessionManager=&SessionManager{}
// 	//SessionManager
// }

// type SessionManager struct {
// 	Session       *ISession
// 	SessionIdName string
// }

// type ISession interface {
// 	Get(name string) interface{}
// 	Set(name string, value interface{})
// 	Remove(name string)
// }

// type InProcSession struct {
// 	innerMap map[string]interface{}
// 	mutex    sync.RWMutex
// }

// func (i *InProcSession) Get(name string) interface{} {
// 	defer i.mutex.RUnlock()
// 	i.mutex.RLock()
// 	v, ok := i.innerMap[name]
// 	if ok {
// 		return v
// 	}
// 	return nil
// }

// func (i *InProcSession) Set(name string, value interface{}) {
// 	i.mutex.Lock()
// 	i.innerMap[name] = value
// 	i.mutex.Unlock()
// }

// func (i *InProcSession) Remove(name string) {
// 	i.mutex.Lock()
// 	delete(i.innerMap[name])
// 	i.mutex.Unlock()
// }
