package handler

import (
	"ai-listen/internal/dto"
	"ai-listen/internal/middleware"
	"ai-listen/internal/pkg/apperror"
	"ai-listen/internal/pkg/response"
	"ai-listen/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) SendCode(c *gin.Context) {
	var req dto.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("参数错误"))
		return
	}

	data, err := h.authService.SendCode(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, data)
}

func (h *AuthHandler) MobileLogin(c *gin.Context) {
	var req dto.MobileLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("参数错误"))
		return
	}

	data, err := h.authService.MobileLogin(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, data)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, err := middleware.CurrentToken(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), token); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{})
}
