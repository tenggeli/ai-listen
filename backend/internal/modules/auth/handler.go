package auth

import (
	"net/http"

	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	handler := &Handler{logger: logger}

	authGroup := group.Group("/auth")
	authGroup.POST("/sms/send", handler.SendSMS)
	authGroup.POST("/login/sms", handler.LoginBySMS)
	authGroup.POST("/login/wechat", handler.LoginByWechat)
	authGroup.POST("/token/refresh", handler.RefreshToken)
	authGroup.POST("/logout", handler.Logout)
}

func (h *Handler) SendSMS(c *gin.Context) {
	var req SendSMSRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "auth", "action": "send_sms", "request": req})
}

func (h *Handler) LoginBySMS(c *gin.Context) {
	var req SMSLoginRequest
	_ = c.ShouldBindJSON(&req)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": gin.H{"module": "auth", "action": "login_sms", "request": req, "todo": "implement token issue"}, "traceId": c.GetHeader("X-Request-Id")})
}

func (h *Handler) LoginByWechat(c *gin.Context) {
	var req WechatLoginRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "auth", "action": "login_wechat", "request": req})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "auth", "action": "refresh_token", "request": req})
}

func (h *Handler) Logout(c *gin.Context) {
	response.Success(c, gin.H{"module": "auth", "action": "logout"})
}
