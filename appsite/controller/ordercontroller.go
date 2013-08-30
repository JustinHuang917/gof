package controller

import (
	"fmt"
	"github.com/JustinHuang917/gof/appsite/models"
	"github.com/JustinHuang917/gof/gofcore"
)

func init() {
	gofcore.RegisterController("order", &OrderController{})
}

type OrderController struct {
	gofcore.ControllerBase
}

func (c OrderController) GetIndex(context *gofcore.HttpContext) (viewResult *gofcore.ViewResult) {
	orders := models.GetAllOrders()
	if orders == nil {
		orders = make([]*models.Order, 0, 0)
	}
	viewResult = gofcore.View(orders, context)
	return
}

func (c OrderController) GetOrder(context *gofcore.HttpContext) (viewResult *gofcore.ViewResult) {
	fmt.Println(context.RoutesData)
	idValue := context.RoutesData.Get("id")
	if id, ok := (idValue).(string); ok {
		order := models.GetOrder(id)
		if order == nil {
			panic("Order not exsited")
		}
		viewResult = gofcore.View(order, context)
	}
	return
}

func (c OrderController) GetCreate(context *gofcore.HttpContext) (viewResult *gofcore.ViewResult) {
	m := &models.Order{}
	viewResult = gofcore.View(m, context)
	return
}

func (c OrderController) PostCreate(context *gofcore.HttpContext, order models.Order) (viewResult *gofcore.ViewResult) {
	id := models.CreateOrder(&order)
	c.RedirectToActionWithRouteData(context, "Order", map[string]string{"id": id})
	return
}
