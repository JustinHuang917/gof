// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/JustinHuang917/gof/goftool"
	"strings"
)

var (
	action      = flag.String("action", "compileview", "compileview")
	viewpath    = flag.String("viewpath", "./view", "view path")
	outviewpath = flag.String("outviewpath", "./view", "out view path")
	others      = flag.String("other", "", "for other args")
)

func main() {
	flag.Parse()
	switch {
	case *action == "compileview":
		otherArgs := make([]string, 0, 0)
		fmt.Println(*others)
		if *others != "" {
			otherArgs = strings.Split(*others, ",")
		}
		err := goftool.Compile(*viewpath, *outviewpath, otherArgs...)
		if err != nil {
			fmt.Println(err)
			return
		}
	case *action == "clearview":
		err := goftool.Clear(*viewpath)
		if err != nil {
			fmt.Println(err)
		}

	}
}
