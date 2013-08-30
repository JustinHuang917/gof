package models

import (
	// "fmt"
	"strconv"
)

type Order struct {
	OrderNo string
	OrderBy string
	Amount  float32
}

var orders []*Order

func CreateOrder(order *Order) string {
	if orders == nil {
		orders = make([]*Order, 0, 10)
	}
	c := len(orders) + 1
	id := strconv.Itoa(c)
	order.OrderNo = id
	orders = append(orders, order)
	return order.OrderNo
}

func GetOrder(id string) *Order {
	for _, order := range orders {
		if order.OrderNo == id {
			return order
		}
	}
	return nil
}

func GetAllOrders() []*Order {
	return orders
}
