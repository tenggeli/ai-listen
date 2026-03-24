package serviceitem

import (
	"ai-listen/backend/internal/store"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	handler := &Handler{logger: logger}
	group.GET("/service-items", handler.List)
}

func (h *Handler) List(c *gin.Context) {
	response.Success(c, gin.H{
		"module":   "service_item",
		"action":   "list",
		"category": c.Query("category"),
		"status":   c.Query("status"),
		"list":     store.Default().ServiceItems(),
	})
}
