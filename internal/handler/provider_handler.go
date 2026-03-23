package handler

import (
	"strconv"

	"ai-listen/internal/pkg/apperror"
	"ai-listen/internal/pkg/response"
	"ai-listen/internal/service"
	"github.com/gin-gonic/gin"
)

type ProviderHandler struct {
	providerService *service.ProviderService
}

func NewProviderHandler(providerService *service.ProviderService) *ProviderHandler {
	return &ProviderHandler{providerService: providerService}
}

func (h *ProviderHandler) List(c *gin.Context) {
	cityID := uint64(0)
	if cityIDStr := c.Query("cityId"); cityIDStr != "" {
		parsed, err := strconv.ParseUint(cityIDStr, 10, 64)
		if err != nil {
			c.Error(apperror.BadRequest("cityId 参数错误"))
			return
		}
		cityID = parsed
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		parsed, err := strconv.Atoi(pageStr)
		if err != nil {
			c.Error(apperror.BadRequest("page 参数错误"))
			return
		}
		page = parsed
	}

	pageSize := 10
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		parsed, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			c.Error(apperror.BadRequest("pageSize 参数错误"))
			return
		}
		pageSize = parsed
	}

	data, err := h.providerService.List(c.Request.Context(), cityID, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, data)
}

func (h *ProviderHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperror.BadRequest("id 参数错误"))
		return
	}

	data, err := h.providerService.Detail(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, data)
}
