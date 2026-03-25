package admin

import (
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) ListComplaints(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_complaints", "query": httpx.PaginationQuery(c)})
}

func (h *Controller) ComplaintDetail(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "complaint_detail", "complaintId": c.Param("complaintId")})
}

func (h *Controller) ResolveComplaint(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "resolve_complaint", "complaintId": c.Param("complaintId")})
}

func (h *Controller) ListRiskEvents(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_risk_events", "query": httpx.PaginationQuery(c)})
}
