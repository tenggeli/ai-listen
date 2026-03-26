package review

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	controller := &Controller{logger: logger}

	group.POST("/orders/:orderId/reviews", controller.Create)
	group.GET("/orders/:orderId/reviews", controller.List)
}
