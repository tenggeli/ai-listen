package admin

import (
	"errors"
	"net/http"
	"strconv"

	app "listen/backend/internal/application/service_item_admin"
	domain "listen/backend/internal/domain/service_item_admin"
)

type ServiceItemController struct {
	listUC   app.ListServiceItemsUseCase
	detailUC app.GetServiceItemDetailUseCase
	statusUC app.UpdateServiceItemStatusUseCase
}

func NewServiceItemController(
	listUC app.ListServiceItemsUseCase,
	detailUC app.GetServiceItemDetailUseCase,
	statusUC app.UpdateServiceItemStatusUseCase,
) ServiceItemController {
	return ServiceItemController{listUC: listUC, detailUC: detailUC, statusUC: statusUC}
}

func (c ServiceItemController) HandleList(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	output, err := c.listUC.Execute(r.Context(), app.ListServiceItemsInput{
		ProviderID: r.URL.Query().Get("provider_id"),
		CategoryID: r.URL.Query().Get("category_id"),
		Status:     r.URL.Query().Get("status"),
		Keyword:    r.URL.Query().Get("keyword"),
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		writeServiceItemError(w, err)
		return
	}

	items := make([]map[string]any, 0, len(output.Items))
	for _, item := range output.Items {
		items = append(items, map[string]any{
			"id":             item.ID,
			"provider_id":    item.ProviderID,
			"provider_name":  item.ProviderName,
			"category_id":    item.CategoryID,
			"title":          item.Title,
			"description":    item.Description,
			"price_amount":   item.PriceAmount,
			"price_unit":     item.PriceUnit,
			"support_online": item.SupportOnline,
			"sort_order":     item.SortOrder,
			"service_status": item.Status,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": output.Total})
}

func (c ServiceItemController) HandleDetail(w http.ResponseWriter, r *http.Request) {
	serviceItemID := r.PathValue("id")
	if serviceItemID == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid service item id")
		return
	}
	output, err := c.detailUC.Execute(r.Context(), app.GetServiceItemDetailInput{ServiceItemID: serviceItemID})
	if err != nil {
		writeServiceItemError(w, err)
		return
	}
	item := output.Item
	writeJSON(w, http.StatusOK, map[string]any{
		"id":             item.ID,
		"provider_id":    item.ProviderID,
		"provider_name":  item.ProviderName,
		"category_id":    item.CategoryID,
		"title":          item.Title,
		"description":    item.Description,
		"price_amount":   item.PriceAmount,
		"price_unit":     item.PriceUnit,
		"support_online": item.SupportOnline,
		"sort_order":     item.SortOrder,
		"service_status": item.Status,
	})
}

func (c ServiceItemController) HandleActivate(w http.ResponseWriter, r *http.Request) {
	c.handleAction(w, r, "activate")
}

func (c ServiceItemController) HandleDeactivate(w http.ResponseWriter, r *http.Request) {
	c.handleAction(w, r, "deactivate")
}

func (c ServiceItemController) handleAction(w http.ResponseWriter, r *http.Request, action string) {
	serviceItemID := r.PathValue("id")
	if serviceItemID == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid service item id")
		return
	}
	output, err := c.statusUC.Execute(r.Context(), app.UpdateServiceItemStatusInput{
		ServiceItemID: serviceItemID,
		Action:        action,
		Operator:      currentAdminID(r),
	})
	if err != nil {
		writeServiceItemError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ServiceItemStatusActionResponseDTO{
		ID:            output.Item.ID,
		ServiceStatus: output.Item.Status,
	})
}

func writeServiceItemError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrServiceItemNotFound):
		writeJSONError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrPersistenceUnavailableInMemory):
		writeJSONError(w, http.StatusNotImplemented, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal error")
	}
}
