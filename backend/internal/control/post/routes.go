package post

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	controller := &Controller{logger: logger}

	group.POST("/posts", controller.Create)
	group.GET("/posts", controller.List)
	group.GET("/posts/:postId", controller.Detail)
	group.POST("/posts/:postId/comments", controller.Comment)
	group.POST("/posts/:postId/likes", controller.Like)
	group.DELETE("/posts/:postId/likes", controller.Unlike)
}
