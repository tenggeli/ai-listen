package me

import (
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	handler := &Controller{logger: logger}
	group.GET("/me/home", handler.Home)
	group.GET("/me/messages", handler.Messages)
	group.GET("/me/settings", handler.Settings)
	group.PUT("/me/settings", handler.UpdateSettings)
}

func (h *Controller) Home(c *gin.Context) {
	response.Success(c, gin.H{"module": "me", "action": "home"})
}

func (h *Controller) Messages(c *gin.Context) {
	response.Success(c, gin.H{"module": "me", "action": "messages"})
}

func (h *Controller) Settings(c *gin.Context) {
	response.Success(c, gin.H{"module": "me", "action": "settings"})
}

func (h *Controller) UpdateSettings(c *gin.Context) {
	response.Success(c, gin.H{"module": "me", "action": "update_settings"})
}
