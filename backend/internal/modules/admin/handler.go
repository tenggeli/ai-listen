package admin

import (
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	handler := &Handler{logger: logger}

	authGroup := group.Group("/auth")
	authGroup.POST("/login", handler.Login)
	authGroup.POST("/logout", handler.Logout)
	authGroup.GET("/me", handler.Me)

	dashboardGroup := group.Group("/dashboard")
	dashboardGroup.GET("/overview", handler.DashboardOverview)

	group.GET("/users", handler.ListUsers)
	group.GET("/users/:userId", handler.UserDetail)
	group.PUT("/users/:userId/status", handler.UpdateUserStatus)

	group.GET("/providers", handler.ListProviders)
	group.GET("/providers/:providerId", handler.ProviderDetail)
	group.POST("/providers/:providerId/approve", handler.ApproveProvider)
	group.POST("/providers/:providerId/reject", handler.RejectProvider)
	group.PUT("/providers/:providerId/status", handler.UpdateProviderStatus)

	group.GET("/service-items", handler.ListServiceItems)
	group.POST("/service-items", handler.CreateServiceItem)
	group.PUT("/service-items/:serviceItemId", handler.UpdateServiceItem)
	group.DELETE("/service-items/:serviceItemId", handler.DeleteServiceItem)

	group.GET("/orders", handler.ListOrders)
	group.GET("/orders/:orderId", handler.OrderDetail)
	group.POST("/orders/:orderId/manual-complete", handler.ManualCompleteOrder)
	group.POST("/orders/:orderId/refund", handler.RefundOrder)

	group.GET("/withdraws", handler.ListWithdraws)
	group.POST("/withdraws/:withdrawId/approve", handler.ApproveWithdraw)
	group.POST("/withdraws/:withdrawId/reject", handler.RejectWithdraw)
	group.GET("/finance/reports", handler.FinanceReports)

	group.GET("/posts", handler.ListPosts)
	group.POST("/posts/:postId/hide", handler.HidePost)
	group.GET("/audio", handler.ListAudio)
	group.POST("/audio/:audioId/off-shelf", handler.OffShelfAudio)

	group.GET("/complaints", handler.ListComplaints)
	group.GET("/complaints/:complaintId", handler.ComplaintDetail)
	group.POST("/complaints/:complaintId/resolve", handler.ResolveComplaint)
	group.GET("/risk-events", handler.ListRiskEvents)

	group.GET("/configs", handler.ListConfigs)
	group.PUT("/configs/:configKey", handler.UpdateConfig)
}

func (h *Handler) Login(c *gin.Context) { httpx.NotImplemented(c, "admin.auth.login") }
func (h *Handler) Logout(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "logout"})
}
func (h *Handler) Me(c *gin.Context) { response.Success(c, gin.H{"module": "admin", "action": "me"}) }
func (h *Handler) DashboardOverview(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "dashboard_overview"})
}
func (h *Handler) ListUsers(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_users", "query": httpx.PaginationQuery(c)})
}
func (h *Handler) UserDetail(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "user_detail", "userId": c.Param("userId")})
}
func (h *Handler) UpdateUserStatus(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "update_user_status", "userId": c.Param("userId")})
}
func (h *Handler) ListProviders(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_providers", "query": httpx.PaginationQuery(c)})
}
func (h *Handler) ProviderDetail(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "provider_detail", "providerId": c.Param("providerId")})
}
func (h *Handler) ApproveProvider(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "approve_provider", "providerId": c.Param("providerId")})
}
func (h *Handler) RejectProvider(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "reject_provider", "providerId": c.Param("providerId")})
}
func (h *Handler) UpdateProviderStatus(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "update_provider_status", "providerId": c.Param("providerId")})
}
func (h *Handler) ListServiceItems(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_service_items"})
}
func (h *Handler) CreateServiceItem(c *gin.Context) {
	httpx.NotImplemented(c, "admin.create_service_item")
}
func (h *Handler) UpdateServiceItem(c *gin.Context) {
	httpx.NotImplemented(c, "admin.update_service_item")
}
func (h *Handler) DeleteServiceItem(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "delete_service_item", "serviceItemId": c.Param("serviceItemId")})
}
func (h *Handler) ListOrders(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_orders", "query": httpx.PaginationQuery(c)})
}
func (h *Handler) OrderDetail(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "order_detail", "orderId": c.Param("orderId")})
}
func (h *Handler) ManualCompleteOrder(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "manual_complete_order", "orderId": c.Param("orderId")})
}
func (h *Handler) RefundOrder(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "refund_order", "orderId": c.Param("orderId")})
}
func (h *Handler) ListWithdraws(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_withdraws", "query": httpx.PaginationQuery(c)})
}
func (h *Handler) ApproveWithdraw(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "approve_withdraw", "withdrawId": c.Param("withdrawId")})
}
func (h *Handler) RejectWithdraw(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "reject_withdraw", "withdrawId": c.Param("withdrawId")})
}
func (h *Handler) FinanceReports(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "finance_reports"})
}
func (h *Handler) ListPosts(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_posts", "query": httpx.PaginationQuery(c)})
}
func (h *Handler) HidePost(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "hide_post", "postId": c.Param("postId")})
}
func (h *Handler) ListAudio(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_audio", "query": httpx.PaginationQuery(c)})
}
func (h *Handler) OffShelfAudio(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "off_shelf_audio", "audioId": c.Param("audioId")})
}
func (h *Handler) ListComplaints(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_complaints", "query": httpx.PaginationQuery(c)})
}
func (h *Handler) ComplaintDetail(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "complaint_detail", "complaintId": c.Param("complaintId")})
}
func (h *Handler) ResolveComplaint(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "resolve_complaint", "complaintId": c.Param("complaintId")})
}
func (h *Handler) ListRiskEvents(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_risk_events", "query": httpx.PaginationQuery(c)})
}
func (h *Handler) ListConfigs(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_configs"})
}
func (h *Handler) UpdateConfig(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "update_config", "configKey": c.Param("configKey")})
}
