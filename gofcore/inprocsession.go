// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofcore

import (
	// "fmt"
	"sync"
)

type InProcSession struct {
	innerMap map[string]interface{}
	mutex    sync.RWMutex
}

func NewInProcSession() *InProcSession {
	s := &InProcSession{}
	s.innerMap = make(map[string]interface{}, 0)
	s.mutex = *new(sync.RWMutex)
	return s
}

func (i InProcSession) Get(name string) interface{} {
	defer i.mutex.RUnlock()
	i.mutex.RLock()
	v, ok := i.innerMap[name]
	if ok {
		return v
	}
	return nil
}

func (i InProcSession) Set(name string, value interface{}) {
	i.mutex.Lock()
	i.innerMap[name] = value
	i.mutex.Unlock()
}

func (i InProcSession) Remove(name string) {
	i.mutex.Lock()
	_, ok1 := i.innerMap[name]
	if ok1 {
		delete(i.innerMap, name)
	}
	i.mutex.Unlock()
}
