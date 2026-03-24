package auth

import (
	"net/http"

	"ai-listen/backend/internal/store"
	"ai-listen/backend/pkg/ecode"
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
	if err := c.ShouldBindJSON(&req); err != nil || req.Mobile == "" {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "mobile is required"})
		return
	}
	code := store.Default().IssueSMSCode(req.Mobile)
	response.Success(c, gin.H{"module": "auth", "action": "send_sms", "request": req, "debugCode": code})
}

func (h *Handler) LoginBySMS(c *gin.Context) {
	var req SMSLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Mobile == "" || req.Code == "" {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "mobile and code are required"})
		return
	}

	user, token, refreshToken, err := store.Default().LoginBySMS(req.Mobile, req.Code)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, ecode.SMSCodeInvalid, gin.H{"reason": err.Error()})
		return
	}

	response.Success(c, gin.H{
		"module":       "auth",
		"action":       "login_sms",
		"user":         user,
		"accessToken":  token,
		"refreshToken": refreshToken,
	})
}

func (h *Handler) LoginByWechat(c *gin.Context) {
	var req WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Code == "" {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "code is required"})
		return
	}

	store.Default().IssueSMSCode("wechat-" + req.Code)
	user, token, refreshToken, err := store.Default().LoginBySMS("wechat-"+req.Code, "123456")
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, ecode.Unauthorized, gin.H{"reason": err.Error()})
		return
	}

	response.Success(c, gin.H{
		"module":       "auth",
		"action":       "login_wechat",
		"user":         user,
		"accessToken":  token,
		"refreshToken": refreshToken,
	})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "refreshToken is required"})
		return
	}

	token, err := store.Default().RefreshToken(req.RefreshToken)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, ecode.Unauthorized, gin.H{"reason": err.Error()})
		return
	}
	response.Success(c, gin.H{"module": "auth", "action": "refresh_token", "accessToken": token})
}

func (h *Handler) Logout(c *gin.Context) {
	response.Success(c, gin.H{"module": "auth", "action": "logout"})
}
