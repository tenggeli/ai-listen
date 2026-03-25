package order

import (
	"strconv"

	"ai-listen/backend/internal/model"
	"ai-listen/backend/internal/support/authctx"
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) List(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "order", "action": "list", "query": httpx.PaginationQuery(c), "list": model.Default().OrdersByUser(user.ID)})
}

func (h *Controller) Detail(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := model.Default().GetOrder(orderID)
	if err != nil {
		httpx.NotImplemented(c, "order.detail")
		return
	}
	response.Success(c, gin.H{"module": "order", "action": "detail", "order": order})
}

func (h *Controller) StatusLogs(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := model.Default().GetOrder(orderID)
	if err != nil {
		httpx.NotImplemented(c, "order.status_logs")
		return
	}
	response.Success(c, gin.H{"module": "order", "action": "status_logs", "list": order.Logs})
}
