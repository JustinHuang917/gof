package main

import (
	"flag"
	"fmt"
	"github.com/justinhuang917/gof/goftool"
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
	// fmt.Println(*action)
	// var arg0 string
	// if flag.NArg() > 0 {
	// 	arg0 = flag.Arg(0)
	// }
	switch {
	case *action == "compileview":
		otherArgs := make([]string, 0, 0)
		fmt.Println(*others)
		if *others != "" {
			otherArgs = strings.Split(*others, ",")
		}
		//	fmt.Println(otherArgs[0])
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