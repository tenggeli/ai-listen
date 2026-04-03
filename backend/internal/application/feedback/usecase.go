package feedback

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

type IDGenerator interface {
	NewID(prefix string) string
}

type SubmitFeedbackInput struct {
	UserID           string
	OrderID          string
	RatingScore      int
	ReviewTags       []string
	ReviewContent    string
	ComplaintReason  string
	ComplaintContent string
}

type SubmitFeedbackOutput struct {
	Item feedbackDomain.OrderFeedback
}

type GetFeedbackInput struct {
	UserID  string
	OrderID string
}

type GetFeedbackOutput struct {
	Item feedbackDomain.OrderFeedback
}

type SubmitOrderFeedbackUseCase struct {
	feedbackRepo feedbackDomain.Repository
	orderRepo    orderDomain.Repository
	idGenerator  IDGenerator
	clock        Clock
}

func NewSubmitOrderFeedbackUseCase(
	feedbackRepo feedbackDomain.Repository,
	orderRepo orderDomain.Repository,
	idGenerator IDGenerator,
	clock Clock,
) SubmitOrderFeedbackUseCase {
	return SubmitOrderFeedbackUseCase{
		feedbackRepo: feedbackRepo,
		orderRepo:    orderRepo,
		idGenerator:  idGenerator,
		clock:        clock,
	}
}

func (u SubmitOrderFeedbackUseCase) Execute(ctx context.Context, input SubmitFeedbackInput) (SubmitFeedbackOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	orderID := strings.TrimSpace(input.OrderID)
	if userID == "" || orderID == "" {
		return SubmitFeedbackOutput{}, feedbackDomain.ErrInvalidInput
	}

	orderItem, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		if err == orderDomain.ErrOrderNotFound {
			return SubmitFeedbackOutput{}, feedbackDomain.ErrFeedbackNotFound
		}
		return SubmitFeedbackOutput{}, err
	}
	if orderItem.UserID != userID {
		return SubmitFeedbackOutput{}, feedbackDomain.ErrFeedbackForbidden
	}
	if orderItem.Status != orderDomain.StatusPaid && orderItem.Status != orderDomain.StatusCompleted {
		return SubmitFeedbackOutput{}, feedbackDomain.ErrOrderNotFinishable
	}

	if _, err := u.feedbackRepo.GetByOrderID(ctx, orderID); err == nil {
		return SubmitFeedbackOutput{}, feedbackDomain.ErrFeedbackSubmitted
	} else if err != feedbackDomain.ErrFeedbackNotFound {
		return SubmitFeedbackOutput{}, err
	}

	item, err := feedbackDomain.NewOrderFeedback(
		u.idGenerator.NewID("fb"),
		orderID,
		userID,
		input.RatingScore,
		input.ReviewTags,
		input.ReviewContent,
		input.ComplaintReason,
		input.ComplaintContent,
		u.clock.Now(),
	)
	if err != nil {
		return SubmitFeedbackOutput{}, err
	}

	hasComplaint := strings.TrimSpace(input.ComplaintReason) != "" || strings.TrimSpace(input.ComplaintContent) != ""
	if hasComplaint {
		if err := orderItem.MarkAfterSaleProcessing(u.clock.Now(), "用户发起投诉"); err != nil {
			return SubmitFeedbackOutput{}, err
		}
		if err := u.orderRepo.Save(ctx, orderItem); err != nil {
			return SubmitFeedbackOutput{}, err
		}
	}

	if err := u.feedbackRepo.Create(ctx, item); err != nil {
		return SubmitFeedbackOutput{}, err
	}
	return SubmitFeedbackOutput{Item: item}, nil
}

type GetOrderFeedbackUseCase struct {
	feedbackRepo feedbackDomain.Repository
	orderRepo    orderDomain.Repository
}

func NewGetOrderFeedbackUseCase(feedbackRepo feedbackDomain.Repository, orderRepo orderDomain.Repository) GetOrderFeedbackUseCase {
	return GetOrderFeedbackUseCase{
		feedbackRepo: feedbackRepo,
		orderRepo:    orderRepo,
	}
}

func (u GetOrderFeedbackUseCase) Execute(ctx context.Context, input GetFeedbackInput) (GetFeedbackOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	orderID := strings.TrimSpace(input.OrderID)
	if userID == "" || orderID == "" {
		return GetFeedbackOutput{}, feedbackDomain.ErrInvalidInput
	}

	orderItem, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		if err == orderDomain.ErrOrderNotFound {
			return GetFeedbackOutput{}, feedbackDomain.ErrFeedbackNotFound
		}
		return GetFeedbackOutput{}, err
	}
	if orderItem.UserID != userID {
		return GetFeedbackOutput{}, feedbackDomain.ErrFeedbackForbidden
	}

	item, err := u.feedbackRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return GetFeedbackOutput{}, err
	}
	return GetFeedbackOutput{Item: item}, nil
}
