package payment

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ai-listen/backend/internal/model"
	"ai-listen/backend/internal/support/authctx"
	"ai-listen/backend/pkg/ecode"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger         *zap.Logger
	callbackSecret string
}

func (h *Controller) Create(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.OrderID == 0 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "orderId is required"})
		return
	}

	payment, order, err := model.Default().CreatePayment(req.OrderID)
	if err != nil {
		h.handlePaymentError(c, err)
		return
	}
	response.Success(c, gin.H{"module": "payment", "action": "create", "payment": payment, "order": order})
}

func (h *Controller) Detail(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	paymentID, err := strconv.ParseUint(c.Param("paymentId"), 10, 64)
	if err != nil || paymentID == 0 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "paymentId is invalid"})
		return
	}

	payment, err := model.Default().GetPayment(paymentID)
	if err != nil {
		h.handlePaymentError(c, err)
		return
	}
	response.Success(c, gin.H{"module": "payment", "action": "detail", "payment": payment})
}

func (h *Controller) Callback(c *gin.Context) {
	var req PaymentCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.PaymentID == 0 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "paymentId is required"})
		return
	}
	if req.Channel == "" || req.NotifyID == "" || req.Sign == "" || req.Timestamp <= 0 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "channel, notifyId, timestamp and sign are required"})
		return
	}
	if abs64(time.Now().Unix()-req.Timestamp) > 300 {
		response.Fail(c, http.StatusUnauthorized, ecode.Unauthorized, gin.H{"reason": "callback timestamp expired"})
		return
	}
	if !verifyCallbackSign(h.callbackSecret, req) {
		response.Fail(c, http.StatusUnauthorized, ecode.Unauthorized, gin.H{"reason": "callback signature invalid"})
		return
	}

	duplicated, err := model.Default().AcquirePaymentCallbackKey(req.Channel, req.NotifyID, req.PaymentID)
	if err != nil {
		h.handlePaymentError(c, err)
		return
	}
	if duplicated {
		payment, order, loadErr := h.loadPaymentOrder(req.PaymentID)
		if loadErr != nil {
			h.handlePaymentError(c, loadErr)
			return
		}
		response.Success(c, gin.H{
			"module":     "payment",
			"action":     "callback",
			"idempotent": true,
			"payment":    payment,
			"order":      order,
		})
		return
	}

	if !strings.EqualFold(req.PayStatus, "SUCCESS") {
		response.Success(c, gin.H{
			"module":    "payment",
			"action":    "callback",
			"accepted":  true,
			"payStatus": req.PayStatus,
			"message":   "ignored non-success callback",
		})
		return
	}

	payment, order, err := model.Default().ConfirmPaymentSuccess(req.PaymentID, req.ThirdTradeNo, req.NotifyRaw)
	if err != nil {
		h.handlePaymentError(c, err)
		return
	}
	response.Success(c, gin.H{"module": "payment", "action": "callback", "payment": payment, "order": order})
}

func (h *Controller) WalletOverview(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "payment", "action": "wallet_overview"})
}

func (h *Controller) PurchaseVIP(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "payment", "action": "purchase_vip"})
}

func (h *Controller) handlePaymentError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, model.ErrOrderNotFound), errors.Is(err, model.ErrPaymentNotFound):
		response.Fail(c, http.StatusNotFound, ecode.NotFound, gin.H{"reason": err.Error()})
	case errors.Is(err, model.ErrInvalidOrderStatus), errors.Is(err, model.ErrInvalidPaymentStatus):
		response.Fail(c, http.StatusConflict, ecode.PaymentStateInvalid, gin.H{"reason": err.Error()})
	case errors.Is(err, model.ErrPaymentCallbackKeyUse):
		response.Fail(c, http.StatusConflict, ecode.Conflict, gin.H{"reason": err.Error()})
	default:
		response.Fail(c, http.StatusInternalServerError, ecode.InternalServerError, gin.H{"reason": err.Error()})
	}
}

func (h *Controller) loadPaymentOrder(paymentID uint64) (*model.Payment, *model.Order, error) {
	payment, err := model.Default().GetPayment(paymentID)
	if err != nil {
		return nil, nil, err
	}
	order, err := model.Default().GetOrder(payment.OrderID)
	if err != nil {
		return nil, nil, err
	}
	return payment, order, nil
}

func abs64(value int64) int64 {
	if value < 0 {
		return -value
	}
	return value
}
