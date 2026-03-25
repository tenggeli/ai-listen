package admin

import (
	"strconv"

	"ai-listen/backend/internal/model"
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) ListProviders(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_providers", "query": httpx.PaginationQuery(c), "list": model.Default().Providers()})
}

func (h *Controller) ProviderDetail(c *gin.Context) {
	providerID, _ := strconv.ParseUint(c.Param("providerId"), 10, 64)
	provider, _ := model.Default().ProviderByID(providerID)
	response.Success(c, gin.H{"module": "admin", "action": "provider_detail", "provider": provider})
}

func (h *Controller) ApproveProvider(c *gin.Context) {
	providerID, _ := strconv.ParseUint(c.Param("providerId"), 10, 64)
	provider, err := model.Default().ApproveProvider(providerID, "approved by admin")
	if err != nil {
		httpx.NotImplemented(c, "admin.approve_provider")
		return
	}
	response.Success(c, gin.H{"module": "admin", "action": "approve_provider", "provider": provider})
}

func (h *Controller) RejectProvider(c *gin.Context) {
	providerID, _ := strconv.ParseUint(c.Param("providerId"), 10, 64)
	provider, err := model.Default().RejectProvider(providerID, "rejected by admin")
	if err != nil {
		httpx.NotImplemented(c, "admin.reject_provider")
		return
	}
	response.Success(c, gin.H{"module": "admin", "action": "reject_provider", "provider": provider})
}

func (h *Controller) UpdateProviderStatus(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "update_provider_status", "providerId": c.Param("providerId")})
}
