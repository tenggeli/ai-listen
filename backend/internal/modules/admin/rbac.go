package admin

import (
	"net/http"

	"ai-listen/backend/internal/store"
	"ai-listen/backend/internal/support/authctx"
	"ai-listen/backend/pkg/ecode"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

const (
	permDashboardOverview = "admin:dashboard:overview"

	permUserList         = "admin:user:list"
	permUserDetail       = "admin:user:detail"
	permUserStatusUpdate = "admin:user:status:update"

	permProviderList         = "admin:provider:list"
	permProviderDetail       = "admin:provider:detail"
	permProviderApprove      = "admin:provider:approve"
	permProviderReject       = "admin:provider:reject"
	permProviderStatusUpdate = "admin:provider:status:update"

	permServiceItemList   = "admin:service_item:list"
	permServiceItemCreate = "admin:service_item:create"
	permServiceItemUpdate = "admin:service_item:update"
	permServiceItemDelete = "admin:service_item:delete"

	permOrderList           = "admin:order:list"
	permOrderDetail         = "admin:order:detail"
	permOrderManualComplete = "admin:order:manual_complete"
	permOrderRefund         = "admin:order:refund"

	permWithdrawList    = "admin:withdraw:list"
	permWithdrawApprove = "admin:withdraw:approve"
	permWithdrawReject  = "admin:withdraw:reject"
	permFinanceReports  = "admin:finance:reports"

	permPostList      = "admin:post:list"
	permPostHide      = "admin:post:hide"
	permAudioList     = "admin:audio:list"
	permAudioOffShelf = "admin:audio:off_shelf"

	permComplaintList    = "admin:complaint:list"
	permComplaintDetail  = "admin:complaint:detail"
	permComplaintResolve = "admin:complaint:resolve"
	permRiskEventList    = "admin:risk_event:list"

	permConfigList   = "admin:config:list"
	permConfigUpdate = "admin:config:update"

	permRBACRoleList       = "admin:rbac:role:list"
	permRBACPermissionList = "admin:rbac:permission:list"
	permRBACUserRoleAssign = "admin:rbac:user_role:assign"
)

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminUser, ok := authctx.CurrentAdmin(c)
		if !ok {
			c.Abort()
			return
		}
		if store.Default().AdminHasPermission(adminUser.ID, permission) {
			c.Next()
			return
		}
		response.Fail(c, http.StatusForbidden, ecode.Forbidden, gin.H{"reason": "admin permission denied", "permission": permission})
		c.Abort()
	}
}

func PermissionsByRoles(roles []string) []string {
	return store.Default().AdminPermissionsByRoles(roles)
}
