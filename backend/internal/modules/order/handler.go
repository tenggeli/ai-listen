package order

import (
	"strconv"

	"ai-listen/backend/internal/store"
	"ai-listen/backend/internal/support/authctx"
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
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req CreateOrderRequest
	_ = c.ShouldBindJSON(&req)
	order, err := store.Default().CreateOrder(user.ID, store.CreateOrderInput{
		ProviderID:      req.ProviderID,
		ServiceItemID:   req.ServiceItemID,
		SceneText:       req.SceneText,
		CityCode:        req.CityCode,
		AddressText:     req.AddressText,
		PlannedStartAt:  req.PlannedStartAt,
		PlannedDuration: req.PlannedDuration,
	})
	if err != nil {
		httpx.NotImplemented(c, "order.create")
		return
	}
	response.Success(c, gin.H{"module": "order", "action": "create", "order": order})
}

func (h *Handler) List(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "order", "action": "list", "query": httpx.PaginationQuery(c), "list": store.Default().OrdersByUser(user.ID)})
}

func (h *Handler) Detail(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := store.Default().GetOrder(orderID)
	if err != nil {
		httpx.NotImplemented(c, "order.detail")
		return
	}
	response.Success(c, gin.H{"module": "order", "action": "detail", "order": order})
}

func (h *Handler) Cancel(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req CancelOrderRequest
	_ = c.ShouldBindJSON(&req)
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := store.Default().CancelOrder(orderID, user.ID, req.Reason)
	if err != nil {
		httpx.NotImplemented(c, "order.cancel")
		return
	}
	response.Success(c, gin.H{"module": "order", "action": "cancel", "order": order})
}

func (h *Handler) Start(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req StartOrderRequest
	_ = c.ShouldBindJSON(&req)
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := store.Default().StartOrder(user.ID, orderID, "user confirmed start service")
	if err != nil {
		httpx.NotImplemented(c, "order.start")
		return
	}
	response.Success(c, gin.H{"module": "order", "action": "start", "order": order, "request": req})
}

func (h *Handler) ConfirmFinish(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := store.Default().ConfirmFinishOrder(user.ID, orderID)
	if err != nil {
		httpx.NotImplemented(c, "order.confirm_finish")
		return
	}
	response.Success(c, gin.H{"module": "order", "action": "confirm_finish", "order": order})
}

func (h *Handler) StatusLogs(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := store.Default().GetOrder(orderID)
	if err != nil {
		httpx.NotImplemented(c, "order.status_logs")
		return
	}
	response.Success(c, gin.H{"module": "order", "action": "status_logs", "list": order.Logs})
}
