package admin

import (
	"net/http"
	"strconv"

	"ai-listen/backend/internal/model"
	"ai-listen/backend/pkg/ecode"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) ListRoles(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_roles", "list": model.Default().AdminRoles()})
}

func (h *Controller) ListPermissions(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_permissions", "list": model.Default().AdminPermissions()})
}

func (h *Controller) AssignUserRoles(c *gin.Context) {
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

	adminUser, err := model.Default().UpdateAdminUserRoles(adminUserID, req.RoleKeys)
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
