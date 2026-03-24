package provider

import (
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	handler := &Handler{logger: logger}

	providersGroup := group.Group("/providers")
	providersGroup.GET("", handler.List)
	providersGroup.GET("/:providerId", handler.Detail)
	providersGroup.GET("/:providerId/posts", handler.Posts)
	providersGroup.POST("/:providerId/favorite", handler.Favorite)
	providersGroup.DELETE("/:providerId/favorite", handler.Unfavorite)
}

func (h *Handler) List(c *gin.Context) {
	response.Success(c, gin.H{
		"module": "provider",
		"action": "list",
		"query": gin.H{
			"serviceItemId": c.Query("serviceItemId"),
			"cityCode":      c.Query("cityCode"),
			"workStatus":    c.Query("workStatus"),
			"minPrice":      c.Query("minPrice"),
			"maxPrice":      c.Query("maxPrice"),
			"sortType":      c.Query("sortType"),
			"pagination":    httpx.PaginationQuery(c),
		},
	})
}

func (h *Handler) Detail(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider", "action": "detail", "providerId": c.Param("providerId")})
}

func (h *Handler) Posts(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider", "action": "posts", "providerId": c.Param("providerId")})
}

func (h *Handler) Favorite(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider", "action": "favorite", "providerId": c.Param("providerId")})
}

func (h *Handler) Unfavorite(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider", "action": "unfavorite", "providerId": c.Param("providerId")})
}
