package admin_order

import (
	"context"
	"strings"
	"time"

	feedbackDomain "listen/backend/internal/domain/feedback"
	orderDomain "listen/backend/internal/domain/order"
)

type Clock interface {
	Now() time.Time
}

type ActionAudit struct {
	ActionID  string
	OrderID   string
	Operator  string
	Reason    string
	UpdatedAt time.Time
	Scope     string
	Action    string
}

type ListOrdersInput struct {
	Status   string
	Keyword  string
	Page     int
	PageSize int
}

type ListOrdersOutput struct {
	Items []orderDomain.Order
	Total int
}

type GetOrderDetailInput struct {
	OrderID string
}

type GetOrderDetailOutput struct {
	Order      orderDomain.Order
	Feedback   *feedbackDomain.OrderFeedback
	ActionLogs []ActionAudit
}

type ActionOrderInput struct {
	OrderID  string
	Action   string
	Operator string
	Reason   string
}

type ActionOrderOutput struct {
	Order orderDomain.Order
	Audit ActionAudit
}

type ComplaintItem struct {
	Order           orderDomain.Order
	Feedback        feedbackDomain.OrderFeedback
	ComplaintStatus string
}

type ListComplaintsInput struct {
	Page     int
	PageSize int
}

type ListComplaintsOutput struct {
	Items []ComplaintItem
	Total int
}

type GetComplaintDetailInput struct {
	OrderID string
}

type GetComplaintDetailOutput struct {
	Item       ComplaintItem
	ActionLogs []ActionAudit
}

type ActionComplaintInput struct {
	OrderID  string
	Action   string
	Operator string
	Reason   string
}

type ActionComplaintOutput struct {
	Item  ComplaintItem
	Audit ActionAudit
}

type UseCase struct {
	orderRepo    orderDomain.Repository
	feedbackRepo feedbackDomain.Repository
	logRepo      ActionLogRepository
	idGenerator  IDGenerator
	clock        Clock
}

type IDGenerator interface {
	NewID(prefix string) string
}

type ActionLogRepository interface {
	Create(ctx context.Context, item ActionAudit) error
	ListByOrderID(ctx context.Context, orderID string) ([]ActionAudit, error)
}

func NewUseCase(
	orderRepo orderDomain.Repository,
	feedbackRepo feedbackDomain.Repository,
	logRepo ActionLogRepository,
	idGenerator IDGenerator,
	clock Clock,
) UseCase {
	return UseCase{
		orderRepo:    orderRepo,
		feedbackRepo: feedbackRepo,
		logRepo:      logRepo,
		idGenerator:  idGenerator,
		clock:        clock,
	}
}

func (u UseCase) ListOrders(ctx context.Context, input ListOrdersInput) (ListOrdersOutput, error) {
	page := input.Page
	if page <= 0 {
		page = 1
	}
	pageSize := input.PageSize
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 10
	}

	items, total, err := u.orderRepo.ListAll(ctx, orderDomain.AdminListQuery{
		Status:   strings.TrimSpace(input.Status),
		Keyword:  strings.TrimSpace(input.Keyword),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return ListOrdersOutput{}, err
	}
	return ListOrdersOutput{Items: items, Total: total}, nil
}

func (u UseCase) GetOrderDetail(ctx context.Context, input GetOrderDetailInput) (GetOrderDetailOutput, error) {
	orderID := strings.TrimSpace(input.OrderID)
	if orderID == "" {
		return GetOrderDetailOutput{}, orderDomain.ErrInvalidInput
	}
	item, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return GetOrderDetailOutput{}, err
	}

	var feedback *feedbackDomain.OrderFeedback
	feedbackItem, feedbackErr := u.feedbackRepo.GetByOrderID(ctx, orderID)
	if feedbackErr == nil {
		feedback = &feedbackItem
	}
	logs, err := u.logRepo.ListByOrderID(ctx, orderID)
	if err != nil {
		return GetOrderDetailOutput{}, err
	}
	return GetOrderDetailOutput{Order: item, Feedback: feedback, ActionLogs: logs}, nil
}

func (u UseCase) ActionOrder(ctx context.Context, input ActionOrderInput) (ActionOrderOutput, error) {
	orderID := strings.TrimSpace(input.OrderID)
	action := strings.TrimSpace(input.Action)
	operator := strings.TrimSpace(input.Operator)
	reason := strings.TrimSpace(input.Reason)
	if orderID == "" || action == "" || operator == "" {
		return ActionOrderOutput{}, orderDomain.ErrInvalidInput
	}

	item, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return ActionOrderOutput{}, err
	}

	switch action {
	case "intervene":
		err = item.MarkAfterSaleProcessing()
	case "close":
		err = item.MarkClosedByAdmin()
	default:
		return ActionOrderOutput{}, orderDomain.ErrInvalidInput
	}
	if err != nil {
		return ActionOrderOutput{}, err
	}

	if err := u.orderRepo.Save(ctx, item); err != nil {
		return ActionOrderOutput{}, err
	}

	now := u.clock.Now()
	audit := ActionAudit{
		ActionID:  u.idGenerator.NewID("oad"),
		OrderID:   orderID,
		Operator:  operator,
		Reason:    reason,
		UpdatedAt: now,
		Scope:     "order",
		Action:    action,
	}
	if err := u.logRepo.Create(ctx, audit); err != nil {
		return ActionOrderOutput{}, err
	}
	return ActionOrderOutput{
		Order: item,
		Audit: audit,
	}, nil
}

func (u UseCase) ListComplaints(ctx context.Context, input ListComplaintsInput) (ListComplaintsOutput, error) {
	page := input.Page
	if page <= 0 {
		page = 1
	}
	pageSize := input.PageSize
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 10
	}

	feedbackItems, total, err := u.feedbackRepo.ListComplaints(ctx, feedbackDomain.ComplaintListQuery{Page: page, PageSize: pageSize})
	if err != nil {
		return ListComplaintsOutput{}, err
	}

	items := make([]ComplaintItem, 0, len(feedbackItems))
	for _, fb := range feedbackItems {
		orderItem, getErr := u.orderRepo.GetByID(ctx, fb.OrderID)
		if getErr != nil {
			continue
		}
		items = append(items, ComplaintItem{
			Order:           orderItem,
			Feedback:        fb,
			ComplaintStatus: complaintStatusFromOrder(orderItem.Status),
		})
	}
	return ListComplaintsOutput{Items: items, Total: total}, nil
}

func (u UseCase) GetComplaintDetail(ctx context.Context, input GetComplaintDetailInput) (GetComplaintDetailOutput, error) {
	orderID := strings.TrimSpace(input.OrderID)
	if orderID == "" {
		return GetComplaintDetailOutput{}, feedbackDomain.ErrInvalidInput
	}

	orderItem, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return GetComplaintDetailOutput{}, err
	}
	feedbackItem, err := u.feedbackRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return GetComplaintDetailOutput{}, err
	}
	if !feedbackItem.HasComplaint {
		return GetComplaintDetailOutput{}, feedbackDomain.ErrFeedbackNotFound
	}

	logs, err := u.logRepo.ListByOrderID(ctx, orderID)
	if err != nil {
		return GetComplaintDetailOutput{}, err
	}

	return GetComplaintDetailOutput{Item: ComplaintItem{
		Order:           orderItem,
		Feedback:        feedbackItem,
		ComplaintStatus: complaintStatusFromOrder(orderItem.Status),
	}, ActionLogs: logs}, nil
}

func (u UseCase) ActionComplaint(ctx context.Context, input ActionComplaintInput) (ActionComplaintOutput, error) {
	orderID := strings.TrimSpace(input.OrderID)
	action := strings.TrimSpace(input.Action)
	operator := strings.TrimSpace(input.Operator)
	reason := strings.TrimSpace(input.Reason)
	if orderID == "" || action == "" || operator == "" {
		return ActionComplaintOutput{}, orderDomain.ErrInvalidInput
	}

	feedbackItem, err := u.feedbackRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return ActionComplaintOutput{}, err
	}
	if !feedbackItem.HasComplaint {
		return ActionComplaintOutput{}, feedbackDomain.ErrFeedbackNotFound
	}

	orderAction := complaintActionToOrderAction(action)
	orderOutput, err := u.ActionOrder(ctx, ActionOrderInput{
		OrderID:  orderID,
		Action:   orderAction,
		Operator: operator,
		Reason:   reason,
	})
	if err != nil {
		return ActionComplaintOutput{}, err
	}

	complaintAudit := ActionAudit{
		ActionID:  u.idGenerator.NewID("cad"),
		OrderID:   orderID,
		Operator:  operator,
		Reason:    reason,
		UpdatedAt: u.clock.Now(),
		Scope:     "complaint",
		Action:    action,
	}
	if err := u.logRepo.Create(ctx, complaintAudit); err != nil {
		return ActionComplaintOutput{}, err
	}

	return ActionComplaintOutput{
		Item: ComplaintItem{
			Order:           orderOutput.Order,
			Feedback:        feedbackItem,
			ComplaintStatus: complaintStatusFromOrder(orderOutput.Order.Status),
		},
		Audit: complaintAudit,
	}, nil
}

func complaintActionToOrderAction(action string) string {
	switch action {
	case "resolve":
		return "close"
	default:
		return "intervene"
	}
}

func complaintStatusFromOrder(status string) string {
	switch status {
	case orderDomain.StatusAfterSale:
		return "processing"
	case orderDomain.StatusClosed:
		return "resolved"
	default:
		return "pending"
	}
}
