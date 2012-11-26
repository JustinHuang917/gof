package view

import (
	"bytes"
	"gofcore/core"
	"models"
	"strings"
)

type V_order_edit struct {
	core.ViewBase
}

func init() {
	core.RegisterViews("/order/edit", &V_order_edit{})
}

func (d *V_order_edit) Render(out *bytes.Buffer, m interface{}, viewBag core.ViewBag, httpContext *core.HttpContext) error {
	if model, ok := m.(*models.User); ok {
		model = model
		_f := func() {
			renderbody := func() {
				d.Writeout(out, "")

				d.Writeout(out, "\n")

				d.Writeout(out, "\n")

				d.Writeout(out, "\n<div>\n	")

				displayDiv := func(s string) {
					d.Writeout(out, "\n	 		<div>")

					d.Writeout(out, s)

					d.Writeout(out, "</div>\n	 	")

				}
				d.Writeout(out, "\n	")

				names := "justinhuang;jordan;mike"
				for _, name := range strings.Split(names, ";") {
					d.Writeout(out, "\n			<p>")

					d.Writeout(out, name)

					d.Writeout(out, "</p>	\n	   ")

				}
				d.Writeout(out, "\n	 ")

				if 1 == 1 {
					d.Writeout(out, "\n	 	<p>")

					displayDiv("xxxxxxxx")
					d.Writeout(out, "</p>\n	 ")

				} else {
					d.Writeout(out, "\n	 	<p>aaaaaaaa</p>\n	 ")

				}
				d.Writeout(out, "\n	 ")

				k := 2
				switch {
				case k == 1:
					d.Writeout(out, "\n	 		<p>k==1</p>\n	 	")

				case k == 2:
					d.Writeout(out, "\n	 		<p>k==2</p>\n	 	")

				}
				d.Writeout(out, "\n\n	 	<div>\n	 		<p>")

				displayDiv("justinhuang")
				d.Writeout(out, "</p>\n	 	</div>\n	 	")

				switch {
				case k == 1:
					d.Writeout(out, "\n	 		<p>k==1</p>\n	 	")

				case k == 2:
					d.Writeout(out, "\n	 		<p>k==2</p>\n	 	")

				}
				d.Writeout(out, "\n\n	")

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
