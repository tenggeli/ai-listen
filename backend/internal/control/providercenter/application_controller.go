package providercenter

import (
	"ai-listen/backend/internal/model"
	"ai-listen/backend/internal/support/authctx"
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) Apply(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req ApplyRequest
	_ = c.ShouldBindJSON(&req)
	provider, err := model.Default().ApplyProvider(user.ID, req.RealName, req.IDCardNo)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.apply")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "apply", "provider": provider})
}

func (h *Controller) ApplyStatus(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	provider, err := model.Default().ProviderByUserID(user.ID)
	if err != nil {
		response.Success(c, gin.H{"module": "provider_center", "action": "apply_status", "provider": nil})
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "apply_status", "provider": provider})
}

func (h *Controller) UpdateProfile(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req UpdateProfileRequest
	_ = c.ShouldBindJSON(&req)
	provider, err := model.Default().UpdateProviderProfile(user.ID, req.DisplayName, req.Intro, req.Tags)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.update_profile")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "update_profile", "provider": provider})
}

func (h *Controller) UpdateServiceItems(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req UpdateServiceItemsRequest
	_ = c.ShouldBindJSON(&req)
	items := make([]model.ProviderServiceItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, model.ProviderServiceItem{
			ServiceItemID: item.ServiceItemID,
			PriceAmount:   item.PriceAmount,
			PriceUnit:     item.PriceUnit,
		})
	}
	provider, err := model.Default().UpdateProviderServiceItems(user.ID, items)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.update_service_items")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "update_service_items", "provider": provider})
}

func (h *Controller) UpdateWorkStatus(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req WorkStatusRequest
	_ = c.ShouldBindJSON(&req)
	provider, err := model.Default().UpdateProviderWorkStatus(user.ID, req.WorkStatus)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.update_work_status")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "update_work_status", "provider": provider})
}
