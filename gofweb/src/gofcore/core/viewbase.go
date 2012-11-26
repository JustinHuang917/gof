package core

import (
	"bytes"
	"errors"
	//"../models"
	"fmt"
	//"net/http"
)

type ViewBase struct {
}

func (d *ViewBase) Writeout(out *bytes.Buffer, content string) {
	//d.Html += content
	out.WriteString(content)

}

func (v *ViewBase) ErrorHandle(msg string) error {
	fmt.Println(msg)
	return errors.New(msg)
}

type ViewResult struct {
	Content *bytes.Buffer
}

func View(v IView, model interface{}, context *HttpContext) (viewResult *ViewResult) {
	viewResult = &ViewResult{Content: new(bytes.Buffer)}
	v.Render(viewResult.Content, model, context.ViewBag, context)
	return
}
