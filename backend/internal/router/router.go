package router

import (
	"net/http"

	"ai-listen/backend/internal/modules/health"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func New(logger *zap.Logger) *gin.Engine {
	gin.SetMode(gin.DebugMode)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    "listen-backend",
			"status":  "ok",
			"message": "listen backend is running",
		})
	})

	api := engine.Group("/api")
	v1 := api.Group("/v1")

	health.RegisterRoutes(v1, logger)

	return engine
}
