package providercenter

import (
	"ai-listen/backend/internal/support/authctx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) CreateAccount(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req AccountRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "provider_center", "action": "create_account", "request": req})
}

func (h *Controller) ListAccounts(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "list_accounts"})
}
