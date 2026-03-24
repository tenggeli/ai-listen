package me

import (
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	handler := &Handler{logger: logger}
	group.GET("/me/home", handler.Home)
	group.GET("/me/messages", handler.Messages)
	group.GET("/me/settings", handler.Settings)
	group.PUT("/me/settings", handler.UpdateSettings)
}

func (h *Handler) Home(c *gin.Context) {
	response.Success(c, gin.H{"module": "me", "action": "home"})
}

func (h *Handler) Messages(c *gin.Context) {
	response.Success(c, gin.H{"module": "me", "action": "messages"})
}

func (h *Handler) Settings(c *gin.Context) {
	response.Success(c, gin.H{"module": "me", "action": "settings"})
}

func (h *Handler) UpdateSettings(c *gin.Context) {
	response.Success(c, gin.H{"module": "me", "action": "update_settings"})
}
