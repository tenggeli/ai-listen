package provider

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	app "listen/backend/internal/application/order"
	domain "listen/backend/internal/domain/order"
)

type OrderController struct {
	listUC   app.ProviderListOrdersUseCase
	detailUC app.ProviderGetOrderUseCase
	actionUC app.ProviderOperateOrderUseCase
}

func NewOrderController(listUC app.ProviderListOrdersUseCase, detailUC app.ProviderGetOrderUseCase, actionUC app.ProviderOperateOrderUseCase) OrderController {
	return OrderController{listUC: listUC, detailUC: detailUC, actionUC: actionUC}
}

func (c OrderController) HandleListOrders(w http.ResponseWriter, r *http.Request) {
	providerID := currentProviderID(r)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	output, err := c.listUC.Execute(r.Context(), app.ProviderListOrdersInput{
		ProviderID: providerID,
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		writeOrderError(w, err)
		return
	}

	items := make([]map[string]any, 0, len(output.Items))
	for _, item := range output.Items {
		items = append(items, buildOrderDTO(item))
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"total": output.Total,
	})
}

func (c OrderController) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	providerID := currentProviderID(r)
	orderID := strings.TrimSpace(r.PathValue("id"))
	if orderID == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	output, err := c.detailUC.Execute(r.Context(), app.ProviderGetOrderInput{ProviderID: providerID, OrderID: orderID})
	if err != nil {
		writeOrderError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, buildOrderDTO(output.Order))
}

func (c OrderController) HandleAccept(w http.ResponseWriter, r *http.Request) {
	c.handleAction(w, r, "accept")
}

func (c OrderController) HandleDepart(w http.ResponseWriter, r *http.Request) {
	c.handleAction(w, r, "depart")
}

func (c OrderController) HandleArrive(w http.ResponseWriter, r *http.Request) {
	c.handleAction(w, r, "arrive")
}

func (c OrderController) HandleStart(w http.ResponseWriter, r *http.Request) {
	c.handleAction(w, r, "start")
}

func (c OrderController) HandleComplete(w http.ResponseWriter, r *http.Request) {
	c.handleAction(w, r, "complete")
}

func (c OrderController) handleAction(w http.ResponseWriter, r *http.Request, action string) {
	providerID := currentProviderID(r)
	orderID := strings.TrimSpace(r.PathValue("id"))

	output, err := c.actionUC.Execute(r.Context(), app.ProviderOperateOrderInput{
		ProviderID: providerID,
		OrderID:    orderID,
		Action:     action,
	})
	if err != nil {
		writeOrderError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, OrderStatusActionResponseDTO{
		ID:                 output.Order.ID,
		Status:             output.Order.Status,
		StatusReason:       domain.StatusReasonByAction(output.Order.Status, output.Order.StatusActionReason),
		StatusActionReason: output.Order.StatusActionReason,
		StatusUpdatedAt:    formatRFC3339Ptr(output.Order.StatusUpdatedAt),
	})
}

func buildOrderDTO(item domain.Order) map[string]any {
	var paidAt *string
	if item.PaidAt != nil {
		formatted := item.PaidAt.Format(time.RFC3339)
		paidAt = &formatted
	}
	return map[string]any{
		"id":                   item.ID,
		"user_id":              item.UserID,
		"provider_id":          item.ProviderID,
		"provider_name":        item.ProviderName,
		"service_item_id":      item.ServiceItemID,
		"service_item_title":   item.ServiceItemTitle,
		"amount":               item.Amount,
		"currency":             item.Currency,
		"status":               item.Status,
		"status_reason":        domain.StatusReasonByAction(item.Status, item.StatusActionReason),
		"status_action_reason": item.StatusActionReason,
		"status_updated_at":    formatRFC3339Ptr(item.StatusUpdatedAt),
		"created_at":           item.CreatedAt.Format(time.RFC3339),
		"paid_at":              paidAt,
	}
}

func formatRFC3339Ptr(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.Format(time.RFC3339)
	return &formatted
}

func writeOrderError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrOrderForbidden):
		writeJSONError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, domain.ErrOrderNotFound):
		writeJSONError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidOrderTransition):
		writeJSONError(w, http.StatusConflict, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal error")
	}
}
