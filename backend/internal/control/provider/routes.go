package provider

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	controller := &Controller{logger: logger}

	providersGroup := group.Group("/providers")
	providersGroup.GET("", controller.List)
	providersGroup.GET("/:providerId", controller.Detail)
	providersGroup.GET("/:providerId/posts", controller.Posts)
	providersGroup.POST("/:providerId/favorite", controller.Favorite)
	providersGroup.DELETE("/:providerId/favorite", controller.Unfavorite)
}
