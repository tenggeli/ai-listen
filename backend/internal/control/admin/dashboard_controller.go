package admin

import (
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) DashboardOverview(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "dashboard_overview"})
}
