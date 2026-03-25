package router

import (
	"net/http"

	"ai-listen/backend/internal/control/admin"
	"ai-listen/backend/internal/control/aimatch"
	"ai-listen/backend/internal/control/audio"
	"ai-listen/backend/internal/control/auth"
	"ai-listen/backend/internal/control/complaint"
	"ai-listen/backend/internal/control/health"
	"ai-listen/backend/internal/control/me"
	"ai-listen/backend/internal/control/order"
	"ai-listen/backend/internal/control/payment"
	"ai-listen/backend/internal/control/post"
	"ai-listen/backend/internal/control/provider"
	"ai-listen/backend/internal/control/providercenter"
	"ai-listen/backend/internal/control/review"
	"ai-listen/backend/internal/control/serviceitem"
	"ai-listen/backend/internal/control/user"
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
	auth.RegisterRoutes(v1, logger)
	user.RegisterRoutes(v1, logger)
	aimatch.RegisterRoutes(v1, logger)
	serviceitem.RegisterRoutes(v1, logger)
	provider.RegisterRoutes(v1, logger)
	providercenter.RegisterRoutes(v1, logger)
	order.RegisterRoutes(v1, logger)
	payment.RegisterRoutes(v1, logger)
	review.RegisterRoutes(v1, logger)
	complaint.RegisterRoutes(v1, logger)
	post.RegisterRoutes(v1, logger)
	audio.RegisterRoutes(v1, logger)
	me.RegisterRoutes(v1, logger)

	adminGroup := v1.Group("/admin")
	admin.RegisterRoutes(adminGroup, logger)

	return engine
}
