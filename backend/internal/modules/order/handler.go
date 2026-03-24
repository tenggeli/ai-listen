package order

import (
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	handler := &Handler{logger: logger}

	ordersGroup := group.Group("/orders")
	ordersGroup.POST("", handler.Create)
	ordersGroup.GET("", handler.List)
	ordersGroup.GET("/:orderId", handler.Detail)
	ordersGroup.POST("/:orderId/cancel", handler.Cancel)
	ordersGroup.POST("/:orderId/start", handler.Start)
	ordersGroup.POST("/:orderId/confirm-finish", handler.ConfirmFinish)
	ordersGroup.GET("/:orderId/status-logs", handler.StatusLogs)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateOrderRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "order", "action": "create", "request": req})
}

func (h *Handler) List(c *gin.Context) {
	response.Success(c, gin.H{"module": "order", "action": "list", "query": httpx.PaginationQuery(c)})
}

func (h *Handler) Detail(c *gin.Context) {
	response.Success(c, gin.H{"module": "order", "action": "detail", "orderId": c.Param("orderId")})
}

func (h *Handler) Cancel(c *gin.Context) {
	var req CancelOrderRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "order", "action": "cancel", "orderId": c.Param("orderId"), "request": req})
}

func (h *Handler) Start(c *gin.Context) {
	var req StartOrderRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "order", "action": "start", "orderId": c.Param("orderId"), "request": req})
}

func (h *Handler) ConfirmFinish(c *gin.Context) {
	response.Success(c, gin.H{"module": "order", "action": "confirm_finish", "orderId": c.Param("orderId")})
}

func (h *Handler) StatusLogs(c *gin.Context) {
	response.Success(c, gin.H{"module": "order", "action": "status_logs", "orderId": c.Param("orderId")})
}
