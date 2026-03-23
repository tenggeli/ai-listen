package middleware

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"ai-listen/internal/pkg/apperror"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	ContextUserIDKey = "current_user_id"
	ContextTokenKey  = "current_token"
)

func AuthRequired(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.Error(apperror.Unauthorized("未登录"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			c.Error(apperror.Unauthorized("无效的认证信息"))
			c.Abort()
			return
		}

		token := strings.TrimSpace(parts[1])
		key := fmt.Sprintf("auth:token:%s", token)
		userIDStr, err := rdb.Get(c.Request.Context(), key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				c.Error(apperror.Unauthorized("登录已失效"))
			} else {
				c.Error(apperror.Internal("鉴权失败", err))
			}
			c.Abort()
			return
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			c.Error(apperror.Internal("鉴权失败", err))
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, userID)
		c.Set(ContextTokenKey, token)
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) (uint64, error) {
	value, ok := c.Get(ContextUserIDKey)
	if !ok {
		return 0, apperror.Unauthorized("未登录")
	}

	userID, ok := value.(uint64)
	if !ok {
		return 0, apperror.Unauthorized("未登录")
	}

	return userID, nil
}

func CurrentToken(c *gin.Context) (string, error) {
	value, ok := c.Get(ContextTokenKey)
	if !ok {
		return "", apperror.Unauthorized("未登录")
	}

	token, ok := value.(string)
	if !ok || token == "" {
		return "", apperror.Unauthorized("未登录")
	}

	return token, nil
}
