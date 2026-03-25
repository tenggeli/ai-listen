package provider

import (
	"strconv"

	"ai-listen/backend/internal/model"
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	handler := &Controller{logger: logger}

	providersGroup := group.Group("/providers")
	providersGroup.GET("", handler.List)
	providersGroup.GET("/:providerId", handler.Detail)
	providersGroup.GET("/:providerId/posts", handler.Posts)
	providersGroup.POST("/:providerId/favorite", handler.Favorite)
	providersGroup.DELETE("/:providerId/favorite", handler.Unfavorite)
}

func (h *Controller) List(c *gin.Context) {
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
		"list": model.Default().Providers(),
	})
}

func (h *Controller) Detail(c *gin.Context) {
	providerID, _ := strconv.ParseUint(c.Param("providerId"), 10, 64)
	provider, err := model.Default().ProviderByID(providerID)
	if err != nil {
		response.Success(c, gin.H{"module": "provider", "action": "detail", "provider": nil})
		return
	}
	response.Success(c, gin.H{"module": "provider", "action": "detail", "provider": provider})
}

func (h *Controller) Posts(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider", "action": "posts", "providerId": c.Param("providerId")})
}

func (h *Controller) Favorite(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider", "action": "favorite", "providerId": c.Param("providerId")})
}

func (h *Controller) Unfavorite(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider", "action": "unfavorite", "providerId": c.Param("providerId")})
}
