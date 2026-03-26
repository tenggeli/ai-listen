package me

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	controller := &Controller{logger: logger}
	group.GET("/me/home", controller.Home)
	group.GET("/me/messages", controller.Messages)
	group.GET("/me/settings", controller.Settings)
	group.PUT("/me/settings", controller.UpdateSettings)
}
