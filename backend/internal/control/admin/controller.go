package admin

import (
	"ai-listen/backend/internal/support/authctx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	controller := &Controller{logger: logger}

	authGroup := group.Group("/auth")
	authGroup.POST("/login", controller.Login)
	authGroup.Use(authctx.AdminRequired())
	authGroup.POST("/logout", controller.Logout)
	authGroup.GET("/me", controller.Me)

	protectedGroup := group.Group("")
	protectedGroup.Use(authctx.AdminRequired())

	dashboardGroup := protectedGroup.Group("/dashboard")
	dashboardGroup.GET("/overview", RequirePermission(permDashboardOverview), controller.DashboardOverview)

	protectedGroup.GET("/users", RequirePermission(permUserList), controller.ListUsers)
	protectedGroup.GET("/users/:userId", RequirePermission(permUserDetail), controller.UserDetail)
	protectedGroup.PUT("/users/:userId/status", RequirePermission(permUserStatusUpdate), controller.UpdateUserStatus)

	protectedGroup.GET("/providers", RequirePermission(permProviderList), controller.ListProviders)
	protectedGroup.GET("/providers/:providerId", RequirePermission(permProviderDetail), controller.ProviderDetail)
	protectedGroup.POST("/providers/:providerId/approve", RequirePermission(permProviderApprove), controller.ApproveProvider)
	protectedGroup.POST("/providers/:providerId/reject", RequirePermission(permProviderReject), controller.RejectProvider)
	protectedGroup.PUT("/providers/:providerId/status", RequirePermission(permProviderStatusUpdate), controller.UpdateProviderStatus)

	protectedGroup.GET("/service-items", RequirePermission(permServiceItemList), controller.ListServiceItems)
	protectedGroup.POST("/service-items", RequirePermission(permServiceItemCreate), controller.CreateServiceItem)
	protectedGroup.PUT("/service-items/:serviceItemId", RequirePermission(permServiceItemUpdate), controller.UpdateServiceItem)
	protectedGroup.DELETE("/service-items/:serviceItemId", RequirePermission(permServiceItemDelete), controller.DeleteServiceItem)

	protectedGroup.GET("/orders", RequirePermission(permOrderList), controller.ListOrders)
	protectedGroup.GET("/orders/:orderId", RequirePermission(permOrderDetail), controller.OrderDetail)
	protectedGroup.POST("/orders/:orderId/manual-complete", RequirePermission(permOrderManualComplete), controller.ManualCompleteOrder)
	protectedGroup.POST("/orders/:orderId/refund", RequirePermission(permOrderRefund), controller.RefundOrder)

	protectedGroup.GET("/withdraws", RequirePermission(permWithdrawList), controller.ListWithdraws)
	protectedGroup.POST("/withdraws/:withdrawId/approve", RequirePermission(permWithdrawApprove), controller.ApproveWithdraw)
	protectedGroup.POST("/withdraws/:withdrawId/reject", RequirePermission(permWithdrawReject), controller.RejectWithdraw)
	protectedGroup.GET("/finance/reports", RequirePermission(permFinanceReports), controller.FinanceReports)

	protectedGroup.GET("/posts", RequirePermission(permPostList), controller.ListPosts)
	protectedGroup.POST("/posts/:postId/hide", RequirePermission(permPostHide), controller.HidePost)
	protectedGroup.GET("/audio", RequirePermission(permAudioList), controller.ListAudio)
	protectedGroup.POST("/audio/:audioId/off-shelf", RequirePermission(permAudioOffShelf), controller.OffShelfAudio)

	protectedGroup.GET("/complaints", RequirePermission(permComplaintList), controller.ListComplaints)
	protectedGroup.GET("/complaints/:complaintId", RequirePermission(permComplaintDetail), controller.ComplaintDetail)
	protectedGroup.POST("/complaints/:complaintId/resolve", RequirePermission(permComplaintResolve), controller.ResolveComplaint)
	protectedGroup.GET("/risk-events", RequirePermission(permRiskEventList), controller.ListRiskEvents)

	protectedGroup.GET("/configs", RequirePermission(permConfigList), controller.ListConfigs)
	protectedGroup.PUT("/configs/:configKey", RequirePermission(permConfigUpdate), controller.UpdateConfig)

	rbacGroup := protectedGroup.Group("/rbac")
	rbacGroup.GET("/roles", RequirePermission(permRBACRoleList), controller.ListRoles)
	rbacGroup.GET("/permissions", RequirePermission(permRBACPermissionList), controller.ListPermissions)
	rbacGroup.PUT("/users/:adminUserId/roles", RequirePermission(permRBACUserRoleAssign), controller.AssignUserRoles)
}
