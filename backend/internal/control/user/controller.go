package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	controller := &Controller{logger: logger}

	usersGroup := group.Group("/users")
	usersGroup.GET("/me", controller.GetMe)
	usersGroup.PUT("/me", controller.UpdateMe)
	usersGroup.GET("/me/profile", controller.GetProfile)
	usersGroup.PUT("/me/profile", controller.UpdateProfile)
	usersGroup.GET("/me/orders", controller.GetOrders)
	usersGroup.GET("/me/favorites", controller.GetFavorites)
}
