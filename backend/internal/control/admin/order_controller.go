package admin

import (
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) ListOrders(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_orders", "query": httpx.PaginationQuery(c)})
}

func (h *Controller) OrderDetail(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "order_detail", "orderId": c.Param("orderId")})
}

func (h *Controller) ManualCompleteOrder(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "manual_complete_order", "orderId": c.Param("orderId")})
}

func (h *Controller) RefundOrder(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "refund_order", "orderId": c.Param("orderId")})
}
