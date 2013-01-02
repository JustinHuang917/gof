package controller

import (
	//"bytes"
	"fmt"
	"github.com/justinhuang917/gof/appsite/models"
	"github.com/justinhuang917/gof/gofcore"
)

func init() {
	gofcore.RegiesterController("home", &HomeController{})
}

type HomeController struct {
	gofcore.ControllerBase
}

func (h *HomeController) Before_Controller_Filter() {

}

func (h *HomeController) After_Controller_Filter() {

}

func (h *HomeController) Before_GetIndex_Filter() {

}

func (h *HomeController) After_GetIndex_Fitler() {

}

func (h HomeController) GetIndex(context *gofcore.HttpContext) (viewResult *gofcore.ViewResult) {
	cid := context.GofSessionId
	v := gofcore.SessionMgr.Get(cid, "username")
	fmt.Println(v)
	if v != nil && v.(string) == "justinhuang" {
		h.RedirectToAction(context, "login")
		return
	}
	v1 := gofcore.GetView(context.RouteName)
	viewResult = gofcore.View(v1.(gofcore.IView), &models.User{"", "", -1}, context)
	return
}

func (h HomeController) PostLogin(context *gofcore.HttpContext, user models.User) (viewResult *gofcore.ViewResult) {
	v := gofcore.GetView(context.RouteName)
	if user.Name == "justinhuang" && user.Password == "123" {
		cid := context.GofSessionId
		gofcore.SessionMgr.Set(cid, "username", "justinhuang")
		viewResult = gofcore.View(v.(gofcore.IView), &models.User{"JustinHuang", "", 100}, context)
	} else {
		h.RedirectToAction(context, "index")
	}
	return
}

func (h HomeController) GetLogin(context *gofcore.HttpContext) (viewResult *gofcore.ViewResult) {
	v := gofcore.GetView(context.RouteName)
	viewResult = gofcore.View(v.(gofcore.IView), &models.User{"JustinHuang", "", 100}, context)
	return
}

func (h HomeController) GetNofound(context *gofcore.HttpContext) (viewResult *gofcore.ViewResult) {
	v := gofcore.GetView(context.RouteName)
	viewResult = gofcore.View(v.(gofcore.IView), gofcore.NullModel, context)
	return
}
