package response

import (
	"net/http"

	"ai-listen/backend/pkg/ecode"
	"github.com/gin-gonic/gin"
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	TraceID string      `json:"traceId"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{
		Code:    ecode.Success.Code,
		Message: ecode.Success.Message,
		Data:    data,
		TraceID: traceID(c),
	})
}

func Fail(c *gin.Context, status int, code ecode.Code, data interface{}) {
	c.JSON(status, Body{
		Code:    code.Code,
		Message: code.Message,
		Data:    data,
		TraceID: traceID(c),
	})
}

func traceID(c *gin.Context) string {
	value := c.GetHeader("X-Request-Id")
	if value == "" {
		return "generated-later"
	}
	return value
}
