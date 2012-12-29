package goftool

var (
	LayoutviewTemplate = `package view

import (
	"github.com/justinhuang917/gof/gofcore"
	"bytes"
	%v
)

type %v struct {
	gofcore.ViewBase
}

func init() {
	gofcore.RegisterViews("%v", &%v{})
}

func (d *%v) Render(out *bytes.Buffer, m interface{},viewBag gofcore.ViewBag, httpContext *gofcore.HttpContext) error {
	if model, ok := m.(*%v); ok {
	model = model
		_f := func() {
			renderbody:=func(){
			 	%v
			 }
			%v	 
		}
		_f()
	 } else {
		errMsg := "The type of model not mtahched"
		return d.ErrorHandle(errMsg)
	}
	return nil
}`

	NoLayoutViewTemplate = `
package view
	
import (
	"github.com/justinhuang917/gof/gofcore"
	"bytes"
	%v
)

type %v struct {
	gofcore.ViewBase
}

func init() {
	gofcore.RegisterViews("%v", &%v{})
}

func (d *%v) Render(out *bytes.Buffer, m interface{},viewBag gofcore.ViewBag, httpContext *gofcore.HttpContext) error {
	if model, ok := m.(*%v); ok {
	model = model
		renderbody := func() {
			%v
		}
		renderbody()
	} else {
		errMsg := "The type of model not matchched"
		return d.ErrorHandle(errMsg)
	}
	return nil
}`
)
