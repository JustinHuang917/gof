package view

import (
	"bytes"
	"github.com/justinhuang917/gof/appsite/models"
	"github.com/justinhuang917/gof/gofcore"
)

type V_home_login struct {
	gofcore.ViewBase
}

func init() {
	gofcore.RegisterViews("/home/login", &V_home_login{})
}

func (d *V_home_login) Render(out *bytes.Buffer, m interface{}, viewBag gofcore.ViewBag, httpContext *gofcore.HttpContext) error {
	if model, ok := m.(*models.User); ok {
		model = model
		_f := func() {
			renderbody := func() {
				d.Writeout(out, "")

				d.Writeout(out, "\n")

				d.Writeout(out, "\n")

				d.Writeout(out, "\n\n")

				displayDiv := func(innertext string) {
					d.Writeout(out, "\n		<div>")

					d.Writeout(out, innertext)

					d.Writeout(out, "</div>\n	 ")

				}
				d.Writeout(out, "")

				displayDiv("Wlecome" + model.Name)
				d.Writeout(out, "")

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
