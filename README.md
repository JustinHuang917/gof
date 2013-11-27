gof
===

GOF: The golang mvc web framework 


###View(Using Razor Template)

define layout page

filename: *.rlayout

``` go
<html>
	<body>
		<div>
			@{renderbody()}
		</div>
	<body>
</html>
```

order browser page:

* fileame format: *.gorazor

* using static type model & layout

``` go

@import "github.com/JustinHuang917/gof/appsite/models"
@model *models.Order
@layout ./appsite/view/html/defaultrazor.rlayout

<div>
	<div>OrderNo:@model.OrderNo</div>
	<div>OrderBy:@model.OrderBy</div>
	<div>Amount:@model.Amount</div>
</div>
```
create order page:
``` go
@import "github.com/JustinHuang917/gof/appsite/models"
@model *models.Order
@layout ./appsite/view/html/defaultrazor.rlayout

<form action="./create" method="post">
	<div>
		OrderBy:<input type="text" id="txtOrderBy" name="OrderBy"/>
		Amount:<input type="text" id="txtAmount" name="Amount"/>
		<input type="Submit" value="Submit"/>
	</div>
</form>
```

###Model
``` go
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
```

###Controller
``` go

func init() {
	gofcore.RegisterController("order", &OrderController{})
}

type OrderController struct {
	gofcore.ControllerBase
}

func (c OrderController) GetOrder(context *gofcore.HttpContext) (viewResult *gofcore.ViewResult) {
	idValue := context.RoutesData.Get("id")
	if id, ok := (idValue).(string); ok {
		order := models.GetOrder(id)
		if order == nil {
			panic("Order not exsited")
		}
		viewResult = c.View(order, context)
	}
	return
}

func (c OrderController) GetCreate(context *gofcore.HttpContext) (viewResult *gofcore.ViewResult) {
	m := &models.Order{}
	viewResult = c.View(m, context)
	return
}

func (c OrderController) PostCreate(context *gofcore.HttpContext, order models.Order) (viewResult *gofcore.ViewResult) {
	id := models.CreateOrder(&order)
	c.RedirectToActionWithRouteData(context, "Order", map[string]string{"id": id})
	return
}

```

### Startup
All view files(*.gorazor)  need to build to go source files,so,just run **./build.sh** to build. 

run **./run.sh** file to startup

access http://localhost:9999/Order/Create

[![xrefs](https://sourcegraph.com/api/repos/github.com/JustinHuang917/gof/badges/xrefs.png)](https://sourcegraph.com/github.com/JustinHuang917/gof)

[![funcs](https://sourcegraph.com/api/repos/github.com/JustinHuang917/gof/badges/funcs.png)](https://sourcegraph.com/github.com/JustinHuang917/gof)

[![top func](https://sourcegraph.com/api/repos/github.com/JustinHuang917/gof/badges/top-func.png)](https://sourcegraph.com/github.com/JustinHuang917/gof)

[![library users](https://sourcegraph.com/api/repos/github.com/JustinHuang917/gof/badges/library-users.png)](https://sourcegraph.com/github.com/JustinHuang917/gof)

[![status](https://sourcegraph.com/api/repos/github.com/JustinHuang917/gof/badges/status.png)](https://sourcegraph.com/github.com/JustinHuang917/gof)

	
