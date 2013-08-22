package controller

import (
	"github.com/JustinHuang917/gof/appsite/models"
	"github.com/JustinHuang917/gof/gofcore"
)

func init() {
	gofcore.RegiesterController("order", &HomeController{})
}

type OrderController struct {
	gofcore.ControllerBase
}

func (h OrderController) GetIndex(context *gofcore.HttpContext) (viewResult *gofcore.ViewResult) {
	v := gofcore.GetView(context.RouteName)
	viewResult = gofcore.View(v.(gofcore.IView), &models.User{"JustinHuang", "", 100, 25}, context)
	return
}

func (h OrderController) PostIndex(context *gofcore.HttpContext) {

}
