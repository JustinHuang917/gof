package view

import (
	"bytes"
	"gofcore/core"
)

type V_order_index struct {
	core.ViewBase
}

func init() {
	core.RegisterViews("/order/index", &V_order_index{})
}

func (d *V_order_index) Render(out *bytes.Buffer, m interface{}, viewBag core.ViewBag, httpContext *core.HttpContext) error {
	if model, ok := m.(*core.NilModel); ok {
		model = model
		renderbody := func() {
			d.Writeout(out, "")

		}
		renderbody()
	} else {
		errMsg := "The type of model not matchched"
		return d.ErrorHandle(errMsg)
	}
	return nil
}
