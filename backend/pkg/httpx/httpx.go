package httpx

import (
	"net/http"

	"ai-listen/backend/pkg/ecode"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func NotImplemented(c *gin.Context, action string) {
	response.Fail(c, http.StatusNotImplemented, ecode.InternalServerError, gin.H{
		"action": action,
		"todo":   "handler skeleton initialized",
	})
}

func PaginationQuery(c *gin.Context) gin.H {
	return gin.H{
		"page":     c.DefaultQuery("page", "1"),
		"pageSize": c.DefaultQuery("pageSize", "20"),
	}
}
