package controller

import (
	//"bytes"
	"gofcore/core"
	//"fmt"
	"models"
)

func init() {
	core.RegiesterController("home", &HomeController{})
}

type HomeController struct {
	core.ControllerBase
}

func (h HomeController) GetIndex(context *core.HttpContext) (viewResult *core.ViewResult) {
	v := core.GetView(context.RouteName)
	viewResult = core.View(v.(core.IView), &models.User{"JustinHuang", 100}, context)
	return
}

func (h HomeController) PostIndex(context *core.HttpContext) {

}

func (h HomeController) GetNofound(context *core.HttpContext) (viewResult *core.ViewResult) {
	v := core.GetView(context.RouteName)
	viewResult = core.View(v.(core.IView), core.NullModel, context)
	return
}
