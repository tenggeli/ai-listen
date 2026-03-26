package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

func NewHandler(logger *zap.Logger) *Controller {
	return &Controller{logger: logger}
}

func (h *Controller) Health(c *gin.Context) {
	h.logger.Debug("health check")

	c.JSON(http.StatusOK, gin.H{
		"service": "listen-backend",
		"status":  "ok",
	})
}
