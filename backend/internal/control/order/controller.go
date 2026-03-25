package order

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	controller := &Controller{logger: logger}

	ordersGroup := group.Group("/orders")
	ordersGroup.POST("", controller.Create)
	ordersGroup.GET("", controller.List)
	ordersGroup.GET("/:orderId", controller.Detail)
	ordersGroup.POST("/:orderId/cancel", controller.Cancel)
	ordersGroup.POST("/:orderId/start", controller.Start)
	ordersGroup.POST("/:orderId/confirm-finish", controller.ConfirmFinish)
	ordersGroup.GET("/:orderId/status-logs", controller.StatusLogs)
}
