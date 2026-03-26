package complaint

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	controller := &Controller{logger: logger}

	group.POST("/orders/:orderId/complaints", controller.Create)
	group.GET("/complaints/:complaintId", controller.Detail)
}
