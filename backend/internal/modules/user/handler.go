package user

import (
	"ai-listen/backend/internal/store"
	"ai-listen/backend/internal/support/authctx"
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
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "user", "action": "get_me", "user": user})
}

func (h *Handler) UpdateMe(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req UpdateMeRequest
	_ = c.ShouldBindJSON(&req)
	updated, err := store.Default().UpdateUser(user.ID, func(u *store.User) {
		if req.Nickname != "" {
			u.Nickname = req.Nickname
		}
		if req.Avatar != "" {
			u.Avatar = req.Avatar
		}
		if req.Gender != 0 {
			u.Gender = req.Gender
		}
		if req.Birthday != "" {
			u.Birthday = req.Birthday
		}
		if req.CityCode != "" {
			u.CityCode = req.CityCode
		}
	})
	if err != nil {
		response.Success(c, gin.H{"module": "user", "action": "update_me", "error": err.Error()})
		return
	}
	response.Success(c, gin.H{"module": "user", "action": "update_me", "user": updated})
}

func (h *Handler) GetProfile(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "user", "action": "get_profile", "userId": user.ID})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req UpdateProfileRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "user", "action": "update_profile", "request": req})
}

func (h *Handler) GetOrders(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "user", "action": "get_orders", "query": httpx.PaginationQuery(c), "list": store.Default().OrdersByUser(user.ID)})
}

func (h *Handler) GetFavorites(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "user", "action": "get_favorites", "query": httpx.PaginationQuery(c)})
}
