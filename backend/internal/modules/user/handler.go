package user

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

	usersGroup := group.Group("/users")
	usersGroup.GET("/me", handler.GetMe)
	usersGroup.PUT("/me", handler.UpdateMe)
	usersGroup.GET("/me/profile", handler.GetProfile)
	usersGroup.PUT("/me/profile", handler.UpdateProfile)
	usersGroup.GET("/me/orders", handler.GetOrders)
	usersGroup.GET("/me/favorites", handler.GetFavorites)
}

func (h *Handler) GetMe(c *gin.Context) {
	response.Success(c, gin.H{"module": "user", "action": "get_me"})
}

func (h *Handler) UpdateMe(c *gin.Context) {
	var req UpdateMeRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "user", "action": "update_me", "request": req})
}

func (h *Handler) GetProfile(c *gin.Context) {
	response.Success(c, gin.H{"module": "user", "action": "get_profile"})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "user", "action": "update_profile", "request": req})
}

func (h *Handler) GetOrders(c *gin.Context) {
	response.Success(c, gin.H{"module": "user", "action": "get_orders", "query": httpx.PaginationQuery(c)})
}

func (h *Handler) GetFavorites(c *gin.Context) {
	response.Success(c, gin.H{"module": "user", "action": "get_favorites", "query": httpx.PaginationQuery(c)})
}
