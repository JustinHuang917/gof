package view

import (
	"bytes"
	"github.com/justinhuang917/gof/appsite/models"
	"github.com/justinhuang917/gof/gofcore"
)

type V_home_index struct {
	gofcore.ViewBase
}

func init() {
	gofcore.RegisterViews("/home/index", &V_home_index{})
}

func (d *V_home_index) Render(out *bytes.Buffer, m interface{}, viewBag gofcore.ViewBag, httpContext *gofcore.HttpContext) error {
	if model, ok := m.(*models.User); ok {
		model = model
		_f := func() {
			renderbody := func() {
				d.Writeout(out, "")

				d.Writeout(out, "\n")

				d.Writeout(out, "\n")

				d.Writeout(out, "\n\n\n<form action=\"./login\" method=\"post\">\n	<div>\n		<input type=\"text\" id=\"txtName\" name=\"Name\"/>\n		<input type=\"text\" id=\"txtPassword\" name=\"Password\"/>\n		<input type=\"Submit\" value=\"Login\"/>\n	</div>\n</form>\n\n")

			}
			d.Writeout(out, "<html>\n	<body>\n		<div>\n			")

			renderbody()
			d.Writeout(out, "\n		</div>\n	<body>\n</html>\n")

		}
		_f()
	} else {
		errMsg := "The type of model not mtahched"
		return d.ErrorHandle(errMsg)
	}
	return nil
}
