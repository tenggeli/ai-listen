package admin

import (
	"net/http"
	"strconv"

	"ai-listen/backend/internal/store"
	"ai-listen/backend/internal/support/authctx"
	"ai-listen/backend/pkg/ecode"
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
	authGroup.Use(authctx.AdminRequired())
	authGroup.POST("/logout", handler.Logout)
	authGroup.GET("/me", handler.Me)

	protectedGroup := group.Group("")
	protectedGroup.Use(authctx.AdminRequired())

	dashboardGroup := protectedGroup.Group("/dashboard")
	dashboardGroup.GET("/overview", RequirePermission(permDashboardOverview), handler.DashboardOverview)

	protectedGroup.GET("/users", RequirePermission(permUserList), handler.ListUsers)
	protectedGroup.GET("/users/:userId", RequirePermission(permUserDetail), handler.UserDetail)
	protectedGroup.PUT("/users/:userId/status", RequirePermission(permUserStatusUpdate), handler.UpdateUserStatus)

	protectedGroup.GET("/providers", RequirePermission(permProviderList), handler.ListProviders)
	protectedGroup.GET("/providers/:providerId", RequirePermission(permProviderDetail), handler.ProviderDetail)
	protectedGroup.POST("/providers/:providerId/approve", RequirePermission(permProviderApprove), handler.ApproveProvider)
	protectedGroup.POST("/providers/:providerId/reject", RequirePermission(permProviderReject), handler.RejectProvider)
	protectedGroup.PUT("/providers/:providerId/status", RequirePermission(permProviderStatusUpdate), handler.UpdateProviderStatus)

	protectedGroup.GET("/service-items", RequirePermission(permServiceItemList), handler.ListServiceItems)
	protectedGroup.POST("/service-items", RequirePermission(permServiceItemCreate), handler.CreateServiceItem)
	protectedGroup.PUT("/service-items/:serviceItemId", RequirePermission(permServiceItemUpdate), handler.UpdateServiceItem)
	protectedGroup.DELETE("/service-items/:serviceItemId", RequirePermission(permServiceItemDelete), handler.DeleteServiceItem)

	protectedGroup.GET("/orders", RequirePermission(permOrderList), handler.ListOrders)
	protectedGroup.GET("/orders/:orderId", RequirePermission(permOrderDetail), handler.OrderDetail)
	protectedGroup.POST("/orders/:orderId/manual-complete", RequirePermission(permOrderManualComplete), handler.ManualCompleteOrder)
	protectedGroup.POST("/orders/:orderId/refund", RequirePermission(permOrderRefund), handler.RefundOrder)

	protectedGroup.GET("/withdraws", RequirePermission(permWithdrawList), handler.ListWithdraws)
	protectedGroup.POST("/withdraws/:withdrawId/approve", RequirePermission(permWithdrawApprove), handler.ApproveWithdraw)
	protectedGroup.POST("/withdraws/:withdrawId/reject", RequirePermission(permWithdrawReject), handler.RejectWithdraw)
	protectedGroup.GET("/finance/reports", RequirePermission(permFinanceReports), handler.FinanceReports)

	protectedGroup.GET("/posts", RequirePermission(permPostList), handler.ListPosts)
	protectedGroup.POST("/posts/:postId/hide", RequirePermission(permPostHide), handler.HidePost)
	protectedGroup.GET("/audio", RequirePermission(permAudioList), handler.ListAudio)
	protectedGroup.POST("/audio/:audioId/off-shelf", RequirePermission(permAudioOffShelf), handler.OffShelfAudio)

	protectedGroup.GET("/complaints", RequirePermission(permComplaintList), handler.ListComplaints)
	protectedGroup.GET("/complaints/:complaintId", RequirePermission(permComplaintDetail), handler.ComplaintDetail)
	protectedGroup.POST("/complaints/:complaintId/resolve", RequirePermission(permComplaintResolve), handler.ResolveComplaint)
	protectedGroup.GET("/risk-events", RequirePermission(permRiskEventList), handler.ListRiskEvents)

	protectedGroup.GET("/configs", RequirePermission(permConfigList), handler.ListConfigs)
	protectedGroup.PUT("/configs/:configKey", RequirePermission(permConfigUpdate), handler.UpdateConfig)

	rbacGroup := protectedGroup.Group("/rbac")
	rbacGroup.GET("/roles", RequirePermission(permRBACRoleList), handler.ListRoles)
	rbacGroup.GET("/permissions", RequirePermission(permRBACPermissionList), handler.ListPermissions)
	rbacGroup.PUT("/users/:adminUserId/roles", RequirePermission(permRBACUserRoleAssign), handler.AssignUserRoles)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "username and password are required"})
		return
	}
	adminUser, token, err := store.Default().AdminLogin(req.Username, req.Password)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, ecode.Unauthorized, gin.H{"reason": err.Error()})
		return
	}
	response.Success(c, gin.H{
		"module":      "admin",
		"action":      "login",
		"accessToken": token,
		"adminUser":   adminUser,
		"permissions": PermissionsByRoles(adminUser.Roles),
	})
}
func (h *Handler) Logout(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "logout"})
}
func (h *Handler) Me(c *gin.Context) {
	adminUser, ok := authctx.CurrentAdmin(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{
		"module":      "admin",
		"action":      "me",
		"adminUser":   adminUser,
		"permissions": PermissionsByRoles(adminUser.Roles),
	})
}
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
	response.Success(c, gin.H{"module": "admin", "action": "list_providers", "query": httpx.PaginationQuery(c), "list": store.Default().Providers()})
}
func (h *Handler) ProviderDetail(c *gin.Context) {
	providerID, _ := strconv.ParseUint(c.Param("providerId"), 10, 64)
	provider, _ := store.Default().ProviderByID(providerID)
	response.Success(c, gin.H{"module": "admin", "action": "provider_detail", "provider": provider})
}
func (h *Handler) ApproveProvider(c *gin.Context) {
	providerID, _ := strconv.ParseUint(c.Param("providerId"), 10, 64)
	provider, err := store.Default().ApproveProvider(providerID, "approved by admin")
	if err != nil {
		httpx.NotImplemented(c, "admin.approve_provider")
		return
	}
	response.Success(c, gin.H{"module": "admin", "action": "approve_provider", "provider": provider})
}
func (h *Handler) RejectProvider(c *gin.Context) {
	providerID, _ := strconv.ParseUint(c.Param("providerId"), 10, 64)
	provider, err := store.Default().RejectProvider(providerID, "rejected by admin")
	if err != nil {
		httpx.NotImplemented(c, "admin.reject_provider")
		return
	}
	response.Success(c, gin.H{"module": "admin", "action": "reject_provider", "provider": provider})
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

func (h *Handler) ListRoles(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_roles", "list": store.Default().AdminRoles()})
}

func (h *Handler) ListPermissions(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_permissions", "list": store.Default().AdminPermissions()})
}

func (h *Handler) AssignUserRoles(c *gin.Context) {
	adminUserID, err := strconv.ParseUint(c.Param("adminUserId"), 10, 64)
	if err != nil || adminUserID == 0 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "invalid adminUserId"})
		return
	}

	var req UpdateAdminUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "invalid request body"})
		return
	}

	adminUser, err := store.Default().UpdateAdminUserRoles(adminUserID, req.RoleKeys)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": err.Error()})
		return
	}

	response.Success(c, gin.H{
		"module":      "admin",
		"action":      "assign_user_roles",
		"adminUser":   adminUser,
		"permissions": PermissionsByRoles(adminUser.Roles),
	})
}
