package admin

import (
	"net/http"

	"ai-listen/backend/internal/model"
	"ai-listen/backend/internal/support/authctx"
	"ai-listen/backend/pkg/ecode"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "username and password are required"})
		return
	}
	adminUser, token, err := model.Default().AdminLogin(req.Username, req.Password)
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

func (h *Controller) Logout(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "logout"})
}

func (h *Controller) Me(c *gin.Context) {
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
