package aimatch

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	controller := &Controller{logger: logger}

	aiGroup := group.Group("/ai")
	aiGroup.POST("/match", controller.Match)
	aiGroup.GET("/match/:sessionId", controller.GetMatch)
}
