package gofcore

import (
	"bytes"
	"errors"
	"fmt"
)

type ViewBase struct {
}

func (d *ViewBase) Writeout(out *bytes.Buffer, content interface{}) {
	s := fmt.Sprint(content)
	out.WriteString(s)
}

func (v *ViewBase) ErrorHandle(msg string) error {
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
