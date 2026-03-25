package order

import (
	"strconv"

	"ai-listen/backend/internal/model"
	"ai-listen/backend/internal/support/authctx"
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) Create(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req CreateOrderRequest
	_ = c.ShouldBindJSON(&req)
	order, err := model.Default().CreateOrder(user.ID, model.CreateOrderInput{
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

func (h *Controller) Cancel(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req CancelOrderRequest
	_ = c.ShouldBindJSON(&req)
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := model.Default().CancelOrder(orderID, user.ID, req.Reason)
	if err != nil {
		httpx.NotImplemented(c, "order.cancel")
		return
	}
	response.Success(c, gin.H{"module": "order", "action": "cancel", "order": order})
}

func (h *Controller) Start(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req StartOrderRequest
	_ = c.ShouldBindJSON(&req)
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := model.Default().StartOrder(user.ID, orderID, "user confirmed start service")
	if err != nil {
		httpx.NotImplemented(c, "order.start")
		return
	}
	response.Success(c, gin.H{"module": "order", "action": "start", "order": order, "request": req})
}

func (h *Controller) ConfirmFinish(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := model.Default().ConfirmFinishOrder(user.ID, orderID)
	if err != nil {
		httpx.NotImplemented(c, "order.confirm_finish")
		return
	}
	response.Success(c, gin.H{"module": "order", "action": "confirm_finish", "order": order})
}
