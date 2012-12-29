package view

import (
	"bytes"
	"github.com/justinhuang917/gof/gofcore"
)

type V_home_nofound struct {
	gofcore.ViewBase
}

func init() {
	gofcore.RegisterViews("/home/nofound", &V_home_nofound{})
}

func (d *V_home_nofound) Render(out *bytes.Buffer, m interface{}, viewBag gofcore.ViewBag, httpContext *gofcore.HttpContext) error {
	if model, ok := m.(*gofcore.NilModel); ok {
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
