package health

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	controller := NewHandler(logger)
	group.GET("/health", controller.Health)
}
