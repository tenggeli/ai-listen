package user

import (
	"ai-listen/backend/internal/model"
	"ai-listen/backend/internal/support/authctx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) GetMe(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "user", "action": "get_me", "user": user})
}

func (h *Controller) UpdateMe(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}

	var req UpdateMeRequest
	_ = c.ShouldBindJSON(&req)
	updated, err := model.Default().UpdateUser(user.ID, func(u *model.User) {
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

func (h *Controller) GetProfile(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "user", "action": "get_profile", "userId": user.ID})
}

func (h *Controller) UpdateProfile(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req UpdateProfileRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "user", "action": "update_profile", "request": req})
}
