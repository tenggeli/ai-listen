package payment

import (
	"ai-listen/backend/internal/model"
	"ai-listen/backend/internal/support/authctx"
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	handler := &Controller{logger: logger}

	group.POST("/payments", handler.Create)
	group.GET("/payments/:paymentId", handler.Detail)
	group.POST("/payments/callback", handler.Callback)
	group.GET("/wallet/overview", handler.WalletOverview)
	group.POST("/vip/purchase", handler.PurchaseVIP)
}

func (h *Controller) Create(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req CreatePaymentRequest
	_ = c.ShouldBindJSON(&req)
	payment, order, err := model.Default().CreatePayment(req.OrderID)
	if err != nil {
		httpx.NotImplemented(c, "payment.create")
		return
	}
	response.Success(c, gin.H{"module": "payment", "action": "create", "payment": payment, "order": order})
}

func (h *Controller) Detail(c *gin.Context) {
	response.Success(c, gin.H{"module": "payment", "action": "detail", "paymentId": c.Param("paymentId")})
}

func (h *Controller) Callback(c *gin.Context) {
	response.Success(c, gin.H{"module": "payment", "action": "callback"})
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
