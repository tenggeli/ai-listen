package audio

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	controller := &Controller{logger: logger}

	group.GET("/audio/categories", controller.Categories)
	group.GET("/audio", controller.List)
	group.GET("/audio/:audioId", controller.Detail)
	group.POST("/audio/:audioId/play-logs", controller.PlayLog)
	group.POST("/audio/:audioId/favorite", controller.Favorite)
}
