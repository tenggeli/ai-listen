package providercenter

import (
	"strconv"

	"ai-listen/backend/internal/model"
	"ai-listen/backend/internal/support/authctx"
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) GetOrders(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	list, err := model.Default().OrdersByProvider(user.ID)
	if err != nil {
		response.Success(c, gin.H{"module": "provider_center", "action": "get_orders", "list": []any{}})
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "get_orders", "query": httpx.PaginationQuery(c), "list": list})
}

func (h *Controller) AcceptOrder(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := model.Default().ProviderAcceptOrder(user.ID, orderID)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.accept_order")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "accept_order", "order": order})
}

func (h *Controller) DepartOrder(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := model.Default().ProviderDepartOrder(user.ID, orderID)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.depart_order")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "depart_order", "order": order})
}

func (h *Controller) ArriveOrder(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := model.Default().ProviderArriveOrder(user.ID, orderID)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.arrive_order")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "arrive_order", "order": order})
}

func (h *Controller) FinishOrder(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := model.Default().ProviderFinishOrder(user.ID, orderID)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.finish_order")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "finish_order", "order": order})
}
