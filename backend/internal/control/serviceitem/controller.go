package serviceitem

import (
	"ai-listen/backend/internal/model"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

func (h *Controller) List(c *gin.Context) {
	response.Success(c, gin.H{
		"module":   "service_item",
		"action":   "list",
		"category": c.Query("category"),
		"status":   c.Query("status"),
		"list":     model.Default().ServiceItems(),
	})
}
