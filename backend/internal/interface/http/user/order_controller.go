package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	app "listen/backend/internal/application/order"
	domain "listen/backend/internal/domain/order"
)

type OrderController struct {
	createOrderUC app.CreateOrderUseCase
	listOrdersUC  app.ListOrdersUseCase
	getOrderUC    app.GetOrderUseCase
	payOrderUC    app.PayOrderMockSuccessUseCase
}

func NewOrderController(
	createOrderUC app.CreateOrderUseCase,
	listOrdersUC app.ListOrdersUseCase,
	getOrderUC app.GetOrderUseCase,
	payOrderUC app.PayOrderMockSuccessUseCase,
) OrderController {
	return OrderController{
		createOrderUC: createOrderUC,
		listOrdersUC:  listOrdersUC,
		getOrderUC:    getOrderUC,
		payOrderUC:    payOrderUC,
	}
}

func (c OrderController) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := currentUserID(r)

	var body CreateOrderRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := c.createOrderUC.Execute(r.Context(), app.CreateOrderInput{
		UserID:           userID,
		ProviderID:       body.ProviderID,
		ProviderName:     body.ProviderName,
		ServiceItemID:    body.ServiceItemID,
		ServiceItemTitle: body.ServiceItemTitle,
		Amount:           body.Amount,
		Currency:         body.Currency,
	})
	if err != nil {
		writeOrderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildOrderResponse(output.Order))
}

func (c OrderController) HandleListOrders(w http.ResponseWriter, r *http.Request) {
	userID := currentUserID(r)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	output, err := c.listOrdersUC.Execute(r.Context(), app.ListOrdersInput{
		UserID:   userID,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		writeOrderError(w, err)
		return
	}

	items := make([]OrderResponseDTO, 0, len(output.Items))
	for _, item := range output.Items {
		items = append(items, buildOrderResponse(item))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"total": output.Total,
	})
}

func (c OrderController) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	userID := currentUserID(r)
	orderID := strings.TrimSpace(r.PathValue("id"))

	output, err := c.getOrderUC.Execute(r.Context(), app.GetOrderInput{
		UserID:  userID,
		OrderID: orderID,
	})
	if err != nil {
		writeOrderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildOrderResponse(output.Order))
}

func (c OrderController) HandlePayMockSuccess(w http.ResponseWriter, r *http.Request) {
	userID := currentUserID(r)
	orderID := strings.TrimSpace(r.PathValue("id"))

	output, err := c.payOrderUC.Execute(r.Context(), app.PayOrderMockSuccessInput{
		UserID:  userID,
		OrderID: orderID,
	})
	if err != nil {
		writeOrderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildOrderResponse(output.Order))
}

func buildOrderResponse(item domain.Order) OrderResponseDTO {
	var paidAt *string
	var statusUpdatedAt *string
	if item.PaidAt != nil {
		formatted := item.PaidAt.Format(time.RFC3339)
		paidAt = &formatted
	}
	if item.StatusUpdatedAt != nil {
		formatted := item.StatusUpdatedAt.Format(time.RFC3339)
		statusUpdatedAt = &formatted
	}
	return OrderResponseDTO{
		ID:                 item.ID,
		UserID:             item.UserID,
		ProviderID:         item.ProviderID,
		ProviderName:       item.ProviderName,
		ServiceItemID:      item.ServiceItemID,
		ServiceItemTitle:   item.ServiceItemTitle,
		Amount:             item.Amount,
		Currency:           item.Currency,
		Status:             item.Status,
		StatusReason:       domain.StatusReasonByAction(item.Status, item.StatusActionReason),
		StatusActionReason: item.StatusActionReason,
		StatusUpdatedAt:    statusUpdatedAt,
		CreatedAt:          item.CreatedAt.Format(time.RFC3339),
		PaidAt:             paidAt,
	}
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
