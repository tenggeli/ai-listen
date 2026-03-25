package authctx

import (
	"net/http"

	"ai-listen/backend/internal/model"
	"ai-listen/backend/pkg/ecode"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

const (
	userContextKey  = "current_user"
	adminContextKey = "current_admin"
)

func CurrentUser(c *gin.Context) (*model.User, bool) {
	if value, exists := c.Get(userContextKey); exists {
		if user, ok := value.(*model.User); ok && user != nil {
			return user, true
		}
	}

	user, err := model.Default().UserByToken(c.GetHeader("Authorization"))
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, ecode.Unauthorized, gin.H{"reason": err.Error()})
		return nil, false
	}
	c.Set(userContextKey, user)
	return user, true
}

func CurrentAdmin(c *gin.Context) (*model.AdminUser, bool) {
	if value, exists := c.Get(adminContextKey); exists {
		if admin, ok := value.(*model.AdminUser); ok && admin != nil {
			return admin, true
		}
	}

	admin, err := model.Default().AdminByToken(c.GetHeader("Authorization"))
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, ecode.Unauthorized, gin.H{"reason": err.Error()})
		return nil, false
	}
	c.Set(adminContextKey, admin)
	return admin, true
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		admin, ok := CurrentAdmin(c)
		if !ok {
			c.Abort()
			return
		}
		c.Set(adminContextKey, admin)
		c.Next()
	}
}
