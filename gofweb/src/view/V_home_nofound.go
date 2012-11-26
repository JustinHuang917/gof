package view

import (
	"bytes"
	"gofcore/core"
)

type V_home_nofound struct {
	core.ViewBase
}

func init() {
	core.RegisterViews("/home/nofound", &V_home_nofound{})
}

func (d *V_home_nofound) Render(out *bytes.Buffer, m interface{}, viewBag core.ViewBag, httpContext *core.HttpContext) error {
	if model, ok := m.(*core.NilModel); ok {
		model = model
		renderbody := func() {
			d.Writeout(out, "<div>\n	<p>404:No Page Found</p>\n</div>\n")

		}
		renderbody()
	} else {
		errMsg := "The type of model not matchched"
		return d.ErrorHandle(errMsg)
	}
	return nil
}
