package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	app "listen/backend/internal/application/admin_order"
	feedbackDomain "listen/backend/internal/domain/feedback"
	orderDomain "listen/backend/internal/domain/order"
)

type OrderController struct {
	uc app.UseCase
}

func NewOrderController(uc app.UseCase) OrderController {
	return OrderController{uc: uc}
}

func (c OrderController) HandleListOrders(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	output, err := c.uc.ListOrders(r.Context(), app.ListOrdersInput{
		Status:   r.URL.Query().Get("status"),
		Keyword:  r.URL.Query().Get("keyword"),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		writeOrderAdminError(w, err)
		return
	}

	items := make([]map[string]any, 0, len(output.Items))
	for _, item := range output.Items {
		items = append(items, orderToMap(item))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": output.Total})
}

func (c OrderController) HandleGetOrderDetail(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimSpace(r.PathValue("id"))
	output, err := c.uc.GetOrderDetail(r.Context(), app.GetOrderDetailInput{OrderID: orderID})
	if err != nil {
		writeOrderAdminError(w, err)
		return
	}

	data := map[string]any{"order": orderToMap(output.Order)}
	if output.Feedback != nil {
		data["feedback"] = feedbackToMap(*output.Feedback)
	}
	data["action_logs"] = actionLogsToMap(output.ActionLogs)
	writeJSON(w, http.StatusOK, data)
}

func (c OrderController) HandleInterveneOrder(w http.ResponseWriter, r *http.Request) {
	c.handleOrderAction(w, r, "intervene")
}

func (c OrderController) HandleCloseOrder(w http.ResponseWriter, r *http.Request) {
	c.handleOrderAction(w, r, "close")
}

func (c OrderController) HandleListComplaints(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	output, err := c.uc.ListComplaints(r.Context(), app.ListComplaintsInput{Page: page, PageSize: pageSize})
	if err != nil {
		writeOrderAdminError(w, err)
		return
	}

	items := make([]map[string]any, 0, len(output.Items))
	for _, item := range output.Items {
		items = append(items, complaintToMap(item))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": output.Total})
}

func (c OrderController) HandleGetComplaintDetail(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimSpace(r.PathValue("id"))
	output, err := c.uc.GetComplaintDetail(r.Context(), app.GetComplaintDetailInput{OrderID: orderID})
	if err != nil {
		writeOrderAdminError(w, err)
		return
	}
	data := complaintToMap(output.Item)
	data["action_logs"] = actionLogsToMap(output.ActionLogs)
	writeJSON(w, http.StatusOK, data)
}

func (c OrderController) HandleInterveneComplaint(w http.ResponseWriter, r *http.Request) {
	c.handleComplaintAction(w, r, "intervene")
}

func (c OrderController) HandleResolveComplaint(w http.ResponseWriter, r *http.Request) {
	c.handleComplaintAction(w, r, "resolve")
}

func (c OrderController) handleOrderAction(w http.ResponseWriter, r *http.Request, action string) {
	orderID := strings.TrimSpace(r.PathValue("id"))
	var body ReviewActionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil && err.Error() != "EOF" {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := c.uc.ActionOrder(r.Context(), app.ActionOrderInput{
		OrderID:  orderID,
		Action:   action,
		Operator: currentAdminID(r),
		Reason:   body.Reason,
	})
	if err != nil {
		writeOrderAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"order": orderToMap(output.Order),
		"audit": map[string]any{
			"action_id":  output.Audit.ActionID,
			"scope":      output.Audit.Scope,
			"action":     output.Audit.Action,
			"operator":   output.Audit.Operator,
			"reason":     output.Audit.Reason,
			"updated_at": output.Audit.UpdatedAt.Format(timeRFC3339),
		},
	})
}

func (c OrderController) handleComplaintAction(w http.ResponseWriter, r *http.Request, action string) {
	orderID := strings.TrimSpace(r.PathValue("id"))
	var body ReviewActionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil && err.Error() != "EOF" {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := c.uc.ActionComplaint(r.Context(), app.ActionComplaintInput{
		OrderID:  orderID,
		Action:   action,
		Operator: currentAdminID(r),
		Reason:   body.Reason,
	})
	if err != nil {
		writeOrderAdminError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"complaint": complaintToMap(output.Item),
		"audit": map[string]any{
			"action_id":  output.Audit.ActionID,
			"scope":      output.Audit.Scope,
			"action":     output.Audit.Action,
			"operator":   output.Audit.Operator,
			"reason":     output.Audit.Reason,
			"updated_at": output.Audit.UpdatedAt.Format(timeRFC3339),
		},
	})
}

const timeRFC3339 = "2006-01-02T15:04:05Z07:00"

func orderToMap(item orderDomain.Order) map[string]any {
	data := map[string]any{
		"id":                 item.ID,
		"user_id":            item.UserID,
		"provider_id":        item.ProviderID,
		"provider_name":      item.ProviderName,
		"service_item_id":    item.ServiceItemID,
		"service_item_title": item.ServiceItemTitle,
		"amount":             item.Amount,
		"currency":           item.Currency,
		"status":             item.Status,
		"created_at":         item.CreatedAt.Format(timeRFC3339),
	}
	if item.PaidAt != nil {
		data["paid_at"] = item.PaidAt.Format(timeRFC3339)
	}
	return data
}

func feedbackToMap(item feedbackDomain.OrderFeedback) map[string]any {
	return map[string]any{
		"id":                item.ID,
		"order_id":          item.OrderID,
		"user_id":           item.UserID,
		"rating_score":      item.RatingScore,
		"review_tags":       item.ReviewTags,
		"review_content":    item.ReviewContent,
		"has_complaint":     item.HasComplaint,
		"complaint_reason":  item.ComplaintReason,
		"complaint_content": item.ComplaintContent,
		"created_at":        item.CreatedAt.Format(timeRFC3339),
	}
}

func complaintToMap(item app.ComplaintItem) map[string]any {
	return map[string]any{
		"complaint_status": item.ComplaintStatus,
		"order":            orderToMap(item.Order),
		"feedback":         feedbackToMap(item.Feedback),
	}
}

func actionLogsToMap(items []app.ActionAudit) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]any{
			"action_id":  item.ActionID,
			"order_id":   item.OrderID,
			"scope":      item.Scope,
			"action":     item.Action,
			"operator":   item.Operator,
			"reason":     item.Reason,
			"updated_at": item.UpdatedAt.Format(timeRFC3339),
		})
	}
	return result
}

func writeOrderAdminError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, orderDomain.ErrInvalidInput):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, orderDomain.ErrOrderNotFound), errors.Is(err, feedbackDomain.ErrFeedbackNotFound):
		writeJSONError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, orderDomain.ErrInvalidOrderTransition):
		writeJSONError(w, http.StatusConflict, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal error")
	}
}
