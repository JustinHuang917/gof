package core

import (
	"bytes"
)

type ControllerBase struct {
}

func (c *ControllerBase) View(v IView, model interface{}, context *HttpContext) (viewResult *ViewResult) {
	viewResult = &ViewResult{Content: new(bytes.Buffer)}
	v.Render(viewResult.Content, model, context.ViewBag, context)
	return
}
