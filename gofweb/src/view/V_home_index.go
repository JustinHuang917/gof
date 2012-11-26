package view

import (
	"bytes"
	"gofcore/core"
	"models"
	"strconv"
)

type V_home_index struct {
	core.ViewBase
}

func init() {
	core.RegisterViews("/home/index", &V_home_index{})
}

func (d *V_home_index) Render(out *bytes.Buffer, m interface{}, viewBag core.ViewBag, httpContext *core.HttpContext) error {
	if model, ok := m.(*models.User); ok {
		model = model
		_f := func() {
			renderbody := func() {
				d.Writeout(out, "\n")

				d.Writeout(out, "\n")

				d.Writeout(out, "\n")

				d.Writeout(out, "\n<div>\n	")

				displayDiv := func(s string) {
					d.Writeout(out, "\n	 		<div>")

					d.Writeout(out, s)

					d.Writeout(out, "</div>\n	 	")

				}
				d.Writeout(out, "\n\n	 <div>\n	 	<ul>\n	 		<li>UserName:")

				d.Writeout(out, model.Name)

				d.Writeout(out, "</li>\n	 		<li>ID:")

				d.Writeout(out, strconv.Itoa(model.Id))

				d.Writeout(out, "</li>\n	 	<ul>\n	 </div>\n	 ")

				displayDiv("Welcome!!!")
				d.Writeout(out, "\n</div>\n")

			}
			d.Writeout(out, "<html>\n	<body>\n		<div>\n			")

			renderbody()
			d.Writeout(out, "\n		</div>\n	<body>\n</html>")

		}
		_f()
	} else {
		errMsg := "The type of model not mtahched"
		return d.ErrorHandle(errMsg)
	}
	return nil
}
