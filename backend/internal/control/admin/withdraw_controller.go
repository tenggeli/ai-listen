package admin

import (
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) ListWithdraws(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_withdraws", "query": httpx.PaginationQuery(c)})
}

func (h *Controller) ApproveWithdraw(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "approve_withdraw", "withdrawId": c.Param("withdrawId")})
}

func (h *Controller) RejectWithdraw(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "reject_withdraw", "withdrawId": c.Param("withdrawId")})
}

func (h *Controller) FinanceReports(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "finance_reports"})
}
