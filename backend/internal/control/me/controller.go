package me

import (
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
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
