package payment

import (
	"ai-listen/backend/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	cfg := config.Load()
	controller := &Controller{
		logger:         logger,
		callbackSecret: cfg.PaymentCallbackSecret,
	}

	group.POST("/payments", controller.Create)
	group.GET("/payments/:paymentId", controller.Detail)
	group.POST("/payments/callback", controller.Callback)
	group.GET("/wallet/overview", controller.WalletOverview)
	group.POST("/vip/purchase", controller.PurchaseVIP)
}
