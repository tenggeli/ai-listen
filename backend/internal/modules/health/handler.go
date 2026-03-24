package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{logger: logger}
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	handler := NewHandler(logger)
	group.GET("/health", handler.Health)
}

func (h *Handler) Health(c *gin.Context) {
	h.logger.Debug("health check")

	c.JSON(http.StatusOK, gin.H{
		"service": "listen-backend",
		"status":  "ok",
	})
}
