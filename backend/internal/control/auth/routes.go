package auth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	controller := &Controller{logger: logger}

	authGroup := group.Group("/auth")
	authGroup.POST("/sms/send", controller.SendSMS)
	authGroup.POST("/login/sms", controller.LoginBySMS)
	authGroup.POST("/login/wechat", controller.LoginByWechat)
	authGroup.POST("/token/refresh", controller.RefreshToken)
	authGroup.POST("/logout", controller.Logout)
}
