package admin

import (
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) ListUsers(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_users", "query": httpx.PaginationQuery(c)})
}

func (h *Controller) UserDetail(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "user_detail", "userId": c.Param("userId")})
}

func (h *Controller) UpdateUserStatus(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "update_user_status", "userId": c.Param("userId")})
}
