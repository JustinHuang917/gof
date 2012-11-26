package main

import (
	"flag"
	"fmt"
	"goftool"
)

func main() {
	flag.Parse()
	var arg0 string
	if flag.NArg() > 0 {
		arg0 = flag.Arg(0)
	}
	switch {

	case arg0 == "compileview": //Compile View
		viewDir := "./view" //Default view path
		if flag.NArg() > 1 {
			viewDir = flag.Arg(1) //"./view/html/"
		}
		outputdir := viewDir
		if flag.NArg() > 2 {
			outputdir = flag.Arg(2)
		}
		err := goftool.Compile(viewDir, outputdir, flag.Args()[3:]...)
		if err != nil {
			fmt.Println(err)
			return
		}

	case arg0 == "clearview":
		viewDir := "./view" //Default view path
		if flag.NArg() > 1 {
			viewDir = flag.Arg(1) //"./view/html/"
		}
		err := goftool.Clear(viewDir)
		if err != nil {
			fmt.Println(err)
		}

	}
}
