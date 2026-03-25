package providercenter

import (
	"ai-listen/backend/internal/model"
	"ai-listen/backend/internal/support/authctx"
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) CreateInvitation(c *gin.Context) {
	httpx.NotImplemented(c, "provider_center.create_invitation")
}

func (h *Controller) GetWallet(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	provider, err := model.Default().ProviderByUserID(user.ID)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.get_wallet")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "get_wallet", "providerId": provider.ID, "balanceAmount": 0})
}

func (h *Controller) CreateWithdraw(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req WithdrawRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "provider_center", "action": "create_withdraw", "request": req})
}

func (h *Controller) ListWithdraws(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "list_withdraws", "query": httpx.PaginationQuery(c)})
}
