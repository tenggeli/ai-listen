package admin

import (
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) ListServiceItems(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_service_items"})
}

func (h *Controller) CreateServiceItem(c *gin.Context) {
	httpx.NotImplemented(c, "admin.create_service_item")
}

func (h *Controller) UpdateServiceItem(c *gin.Context) {
	httpx.NotImplemented(c, "admin.update_service_item")
}

func (h *Controller) DeleteServiceItem(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "delete_service_item", "serviceItemId": c.Param("serviceItemId")})
}
