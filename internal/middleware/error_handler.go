package middleware

import (
	"fmt"
	"log"

	"ai-listen/internal/pkg/apperror"
	"ai-listen/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() {
			return
		}
		if len(c.Errors) == 0 {
			return
		}

		appErr := apperror.From(c.Errors.Last().Err)
		response.JSON(c, response.HTTPStatusByCode(appErr.Code), appErr.Code, appErr.Message, nil)
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic recovered: %v", r)
				c.Error(apperror.Internal("服务器错误", fmt.Errorf("%v", r)))
				c.Abort()
			}
		}()

		c.Next()
	}
}
