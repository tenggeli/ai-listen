package user

import (
	"ai-listen/backend/internal/model"
	"ai-listen/backend/internal/support/authctx"
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) GetOrders(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "user", "action": "get_orders", "query": httpx.PaginationQuery(c), "list": model.Default().OrdersByUser(user.ID)})
}

func (h *Controller) GetFavorites(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "user", "action": "get_favorites", "query": httpx.PaginationQuery(c)})
}
