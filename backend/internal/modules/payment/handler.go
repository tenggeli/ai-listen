package payment

import (
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

	group.POST("/payments", handler.Create)
	group.GET("/payments/:paymentId", handler.Detail)
	group.POST("/payments/callback", handler.Callback)
	group.GET("/wallet/overview", handler.WalletOverview)
	group.POST("/vip/purchase", handler.PurchaseVIP)
}

func (h *Handler) Create(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req CreatePaymentRequest
	_ = c.ShouldBindJSON(&req)
	payment, order, err := store.Default().CreatePayment(req.OrderID)
	if err != nil {
		httpx.NotImplemented(c, "payment.create")
		return
	}
	response.Success(c, gin.H{"module": "payment", "action": "create", "payment": payment, "order": order})
}

func (h *Handler) Detail(c *gin.Context) {
	response.Success(c, gin.H{"module": "payment", "action": "detail", "paymentId": c.Param("paymentId")})
}

func (h *Handler) Callback(c *gin.Context) {
	response.Success(c, gin.H{"module": "payment", "action": "callback"})
}

func (h *Handler) WalletOverview(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "payment", "action": "wallet_overview"})
}

func (h *Handler) PurchaseVIP(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "payment", "action": "purchase_vip"})
}
