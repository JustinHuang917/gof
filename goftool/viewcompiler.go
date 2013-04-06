// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goftool

import (
	"errors"
	"fmt"
	"github.com/justinhuang917/gof/goftool/parser"
	"github.com/justinhuang917/gof/goftool/parser/html"
	"github.com/justinhuang917/gof/goftool/parser/razor"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	open_tag           = "<%"
	close_tag          = "%>"
	out_tag            = "="
	end_server_tag     = "}"
	import_tag         = "import"
	model_declare_tag  = "model"
	layout_declare_tag = "layout"
	else_tag           = "}else{"
	else_if_tag        = "}elseif{"
	helper_tag         = "helper"
	gohtml_ext         = ".gohtml"
	gorazor_ext        = ".gorazor"
	gorazorlayout_ext  = ".rlayout"
	defaultModel       = "gofcore.NilModel"
	viewDir            = "view/html"
)

var (
	writeout_begin = "d.Writeout(out,"
	writeout_end   = ")\n"
	line           = 0
)

type CompileResult struct {
	ParseResult   *parser.ParseResult
	ViewName      string
	LayoutPath    string
	ModelTypeName string
	RouteName     string
}

type IParser interface {
	Parse() *parser.ParseResult
}

type Compiler struct {
	FileName   string
	FilePath   string
	bytes      []byte
	LineCount  int
	OutputPath string
	Result     *CompileResult
}

func NewCompiler(path, OutputPath string) (*Compiler, error) {
	compiler := &Compiler{}
	fi, err := os.Open(path)
	if err == nil {
		compiler.FileName = fi.Name()
		compiler.FilePath = path
		compiler.OutputPath = OutputPath
		compiler.Result = &CompileResult{}
		compiler.LineCount = 0
		compiler.bytes, err = ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		compiler.Result.ViewName = getViewName(path)
		compiler.Result.LayoutPath = ""
	} else {
		return nil, err
	}
	defer fi.Close()
	return compiler, nil
}

func getParser(path string, content string) IParser {
	extName := filepath.Ext(path)
	if extName == gorazor_ext || extName == gorazorlayout_ext {
		return razor.NewRazorParserEngine(content)
	}
	return html.NewHtmlParserEngine(content)
}

func (c *Compiler) Compile(needOutput bool) (err error) {
	parser := getParser(c.FilePath, string(c.bytes))
	c.Result.RouteName = getRouteName(c.FilePath)
	var html string
	c.Result.ParseResult = parser.Parse()
	if c.Result.ParseResult.LayoutPath != "" {
		c.Result.LayoutPath = c.Result.ParseResult.LayoutPath
	}
	if c.Result.ParseResult.ModelTypeName == "" {
		c.Result.ParseResult.ModelTypeName = defaultModel
	}
	if b := c.Result.ParseResult.IsNeedLayout(); b {
		html, err = c.Result.formatHtmlWithLayout()
	} else {
		html, err = c.Result.formatHtmlNoLayout()
	}
	if err == nil {
		if needOutput {
			err = c.outputWriteFile(html)
		}
	}

	return
}

func (c *CompileResult) formatHtmlNoLayout() (html string, err error) {
	html = fmt.Sprintf(NoLayoutViewTemplate,
		c.ParseResult.Imports,
		c.ViewName,
		c.RouteName,
		c.ViewName,
		c.ViewName,
		c.ParseResult.ModelTypeName,
		c.ParseResult.OutPutContent)
	return
}

func (c *CompileResult) formatHtmlWithLayout() (html string, err error) {
	layoutResult, err2 := compileLayout(c.ParseResult.LayoutPath)
	if err2 == nil {
		layoutBodyResult := c.ParseResult.OutPutContent
		html = fmt.Sprintf(LayoutviewTemplate,
			c.ParseResult.Imports,
			c.ViewName,
			c.RouteName,
			c.ViewName,
			c.ViewName,
			c.ParseResult.ModelTypeName,
			layoutBodyResult,
			layoutResult.ParseResult.OutPutContent)

	} else {
		err = err2
	}
	return
}

func (c *Compiler) outputWriteFile(html string) error {
	os.Create(c.OutputPath)
	outputFile, err2 := os.OpenFile(c.OutputPath, os.O_RDWR, 0666)
	if err2 == nil {
		io.WriteString(outputFile, html)
	}
	defer outputFile.Close()
	return err2
}

var layoutResultCache map[string]*CompileResult

func init() {
	layoutResultCache = make(map[string]*CompileResult)
}

func compileLayout(layoutPath string) (compileResult *CompileResult, err error) {
	if layoutResultCache[layoutPath] == nil {
		c, err1 := NewCompiler(layoutPath, "")
		if err1 != nil {
			err = err1
			return
		}
		err = c.Compile(false)
		compileResult = c.Result
		layoutResultCache[layoutPath] = compileResult
	} else {
		compileResult = layoutResultCache[layoutPath]
	}
	return
}

type visitor struct{}

func (self *visitor) DoCompile(path string, outputDir string, f os.FileInfo, buidingArgs ...string) error {
	if f == nil {
		return nil
	}
	if f.Name() == "view.go" {
		return nil
	}
	if f.IsDir() {
		return nil
	} else if (f.Mode() & os.ModeSymlink) > 0 {
		return nil
	} else {
		var err error
		extName := filepath.Ext(f.Name())
		if extName == gohtml_ext || extName == gorazor_ext {
			OutputPath := filepath.Join(outputDir, getViewName(path)+".go")
			if OutputPath[0] != '.' {
				OutputPath = "./" + OutputPath
			}
			c, err1 := NewCompiler(path, OutputPath)
			if err1 != nil {
				return err1
			}
			fmt.Println("Parsing:", path)
			err = c.Compile(true)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Building:", c.OutputPath)
				br := Building(c.OutputPath, buidingArgs...)
				if br != nil && br.Err != nil {
					//err = genCompileError(path, string(br.Out))
					errMsg := fmt.Sprintf("%s:%s", path, string(br.Out))
					err = errors.New(errMsg)
				}
			}
			if err == nil {
				fmt.Println("Success!!!")
				fmt.Println("------------------------------------------------------")
			}
			return err
		}
	}
	return nil
}

func Compile(dirPath string, outputDir string, buidingArgs ...string) error {
	begin := time.Now()
	v := &visitor{}
	//Clear()
	err := filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
		err1 := v.DoCompile(path, outputDir, f, buidingArgs...)
		return err1
	})
	if err != nil {
		//Clear()
		fmt.Println("-----------------------------------------------")
		fmt.Println("Failed!!!")
	} else {
		fmt.Println("-----------------------------------------------")
		fmt.Println("Build All View Success!!!")
	}
	duration := time.Since(begin)
	fmt.Printf("Time Cost: %vs\n", duration.Seconds())
	return err
}

func Clear(dir string) error {
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return nil
		}
		if f.IsDir() {
			return nil
		} else if (f.Mode() & os.ModeSymlink) > 0 {
			return nil
		} else {
			if ok, _ := filepath.Match("V_*.go", f.Name()); ok {
				err = os.Remove(path)
				return err
			}
		}
		return nil
	})
	return err
}

func getViewName(filepath string) string {
	startIndex := strings.Index(filepath, viewDir)
	startIndex = startIndex + len(viewDir)
	f := filepath[startIndex:]
	viewName := strings.Replace(f, ".", "", -1)
	viewName = strings.Replace(viewName, "/", "_", -1)
	viewName = "V" + strings.Replace(viewName, gohtml_ext[1:], "", -1)
	viewName = strings.Replace(viewName, gorazor_ext[1:], "", -1)
	return viewName
}

func getRouteName(filepath string) string {
	startIndex := strings.Index(filepath, viewDir)
	startIndex = startIndex + len(viewDir)
	f := filepath[startIndex:]
	routeName := strings.Replace(f, ".", "", -1)
	routeName = strings.Replace(routeName, gohtml_ext[1:], "", -1)
	routeName = strings.Replace(routeName, gorazor_ext[1:], "", -1)
	return routeName
}
