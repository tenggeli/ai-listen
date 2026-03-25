package providercenter

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	controller := &Controller{logger: logger}

	centerGroup := group.Group("/provider-center")
	centerGroup.POST("/apply", controller.Apply)
	centerGroup.GET("/apply/status", controller.ApplyStatus)
	centerGroup.PUT("/profile", controller.UpdateProfile)
	centerGroup.PUT("/service-items", controller.UpdateServiceItems)
	centerGroup.PUT("/work-status", controller.UpdateWorkStatus)
	centerGroup.GET("/orders", controller.GetOrders)
	centerGroup.POST("/invitations", controller.CreateInvitation)
	centerGroup.GET("/wallet", controller.GetWallet)
	centerGroup.POST("/accounts", controller.CreateAccount)
	centerGroup.GET("/accounts", controller.ListAccounts)
	centerGroup.POST("/withdraws", controller.CreateWithdraw)
	centerGroup.GET("/withdraws", controller.ListWithdraws)
	centerGroup.POST("/orders/:orderId/accept", controller.AcceptOrder)
	centerGroup.POST("/orders/:orderId/depart", controller.DepartOrder)
	centerGroup.POST("/orders/:orderId/arrive", controller.ArriveOrder)
	centerGroup.POST("/orders/:orderId/finish", controller.FinishOrder)
}
