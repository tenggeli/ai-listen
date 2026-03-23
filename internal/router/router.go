package router

import (
	"net/http"

	"ai-listen/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(
	engine *gin.Engine,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	providerCenterHandler *handler.ProviderCenterHandler,
	homeHandler *handler.HomeHandler,
	providerHandler *handler.ProviderHandler,
	authMiddleware gin.HandlerFunc,
) {
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "api",
			"status":  "ok",
		})
	})

	api := engine.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/send-code", authHandler.SendCode)
			auth.POST("/mobile-login", authHandler.MobileLogin)
			auth.Use(authMiddleware).POST("/logout", authHandler.Logout)
		}

		user := api.Group("/user")
		user.Use(authMiddleware)
		{
			user.GET("/profile", userHandler.Profile)
		}

		providerCenter := api.Group("/provider-center")
		providerCenter.Use(authMiddleware)
		{
			providerCenter.POST("/apply", providerCenterHandler.Apply)
			providerCenter.GET("/audit-status", providerCenterHandler.AuditStatus)
			providerCenter.PUT("/profile", providerCenterHandler.UpdateProfile)
		}

		home := api.Group("/home")
		home.Use(authMiddleware)
		{
			home.POST("/ai-match", homeHandler.AIMatch)
		}

		provider := api.Group("/provider")
		{
			provider.GET("/list", providerHandler.List)
			provider.GET("/detail/:id", providerHandler.Detail)
		}
	}
}
