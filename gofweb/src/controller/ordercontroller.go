package controller

import (
	"gofcore/core"
	"models"
)

func init() {
	core.RegiesterController("order", &HomeController{})
}

type OrderController struct {
	core.ControllerBase
}

func (h OrderController) GetIndex(context *core.HttpContext) (viewResult *core.ViewResult) {
	v := core.GetView(context.RouteName)
	viewResult = core.View(v.(core.IView), &models.User{"JustinHuang", 100}, context)
	return
}

func (h OrderController) PostIndex(context *core.HttpContext) {

}
