// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofcore

import (
	"fmt"
	"github.com/JustinHuang917/gof/gofcore/cfg"
	"strings"
)

const (
	Runtime = 0
	StartUp = 1
)

func Debug(message string, mode int) {
	configMode := strings.ToLower(cfg.AppConfig.DebugMode)
	if configMode == "any" {
		fmt.Println(message)
	} else if mode == 0 && configMode == "runtime" {
		fmt.Println(message)
	} else if mode == 1 && configMode == "startup" {
		fmt.Println(message)
	}
}
