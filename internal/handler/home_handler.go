package handler

import (
	"ai-listen/internal/dto"
	"ai-listen/internal/middleware"
	"ai-listen/internal/pkg/apperror"
	"ai-listen/internal/pkg/response"
	"ai-listen/internal/service"
	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	homeService *service.HomeService
}

func NewHomeHandler(homeService *service.HomeService) *HomeHandler {
	return &HomeHandler{homeService: homeService}
}

func (h *HomeHandler) AIMatch(c *gin.Context) {
	if _, err := middleware.CurrentUserID(c); err != nil {
		c.Error(err)
		return
	}

	var req dto.AIMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("参数错误"))
		return
	}

	data, err := h.homeService.AIMatch(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, data)
}
