package handler

import (
	"ai-listen/internal/dto"
	"ai-listen/internal/middleware"
	"ai-listen/internal/pkg/apperror"
	"ai-listen/internal/pkg/response"
	"ai-listen/internal/service"
	"github.com/gin-gonic/gin"
)

type ProviderCenterHandler struct {
	providerService *service.ProviderCenterService
}

func NewProviderCenterHandler(providerService *service.ProviderCenterService) *ProviderCenterHandler {
	return &ProviderCenterHandler{providerService: providerService}
}

func (h *ProviderCenterHandler) Apply(c *gin.Context) {
	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.ApplyProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("参数错误"))
		return
	}

	data, err := h.providerService.Apply(c.Request.Context(), userID, req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, data)
}

func (h *ProviderCenterHandler) AuditStatus(c *gin.Context) {
	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	data, err := h.providerService.AuditStatus(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, data)
}

func (h *ProviderCenterHandler) UpdateProfile(c *gin.Context) {
	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdateProviderProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("参数错误"))
		return
	}

	data, err := h.providerService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, data)
}
