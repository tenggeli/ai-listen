package authctx

import (
	"net/http"

	"ai-listen/backend/internal/store"
	"ai-listen/backend/pkg/ecode"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context) (*store.User, bool) {
	user, err := store.Default().UserByToken(c.GetHeader("Authorization"))
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, ecode.Unauthorized, gin.H{"reason": err.Error()})
		return nil, false
	}
	return user, true
}
