package model

import (
	"fmt"
	"time"
)

func (s *MemoryStore) CreateOrder(userID uint64, req CreateOrderInput) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	provider, ok := s.providers[req.ProviderID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	if provider.AuditStatus != 2 {
		return nil, ErrProviderAuditRequired
	}
	item, ok := s.serviceItems[req.ServiceItemID]
	if !ok {
		return nil, ErrServiceItemNotFound
	}

	now := time.Now()
	order := &Order{
		ID:              s.nextOrderID,
		OrderNo:         fmt.Sprintf("L%014d", s.nextOrderID),
		UserID:          userID,
		ProviderID:      req.ProviderID,
		ServiceItemID:   req.ServiceItemID,
		SceneText:       req.SceneText,
		CityCode:        req.CityCode,
		AddressText:     req.AddressText,
		PlannedStartAt:  req.PlannedStartAt,
		PlannedDuration: req.PlannedDuration,
		Status:          OrderStatusPendingPayment,
		PayAmount:       item.MinPrice,
		CreatedAt:       now,
		UpdatedAt:       now,
		Logs: []OrderLog{{
			FromStatus:   0,
			ToStatus:     OrderStatusPendingPayment,
			OperatorRole: "user",
			OperatorID:   userID,
			Remark:       "order created",
			CreatedAt:    now,
		}},
	}
	s.orders[order.ID] = order
	s.nextOrderID++
	return cloneOrder(order), nil
}

func (s *MemoryStore) OrdersByUser(userID uint64) []*Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Order
	for _, order := range s.orders {
		if order.UserID == userID {
			result = append(result, cloneOrder(order))
		}
	}
	return result
}

func (s *MemoryStore) OrdersByProvider(userID uint64) ([]*Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	provider, ok := s.providerByUID[userID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	var result []*Order
	for _, order := range s.orders {
		if order.ProviderID == provider.ID {
			result = append(result, cloneOrder(order))
		}
	}
	return result, nil
}

func (s *MemoryStore) GetOrder(orderID uint64) (*Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	return cloneOrder(order), nil
}

func (s *MemoryStore) CreatePayment(orderID uint64) (*Payment, *Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderID]
	if !ok {
		return nil, nil, ErrOrderNotFound
	}
	if order.Status != OrderStatusPendingPayment {
		return nil, nil, ErrInvalidOrderStatus
	}

	now := time.Now()
	payment := &Payment{
		ID:        s.nextPaymentID,
		OrderID:   orderID,
		PayAmount: order.PayAmount,
		Status:    20,
		CreatedAt: now,
		PaidAt:    now,
	}
	s.payments[payment.ID] = payment
	s.nextPaymentID++
	s.transitionOrder(order, OrderStatusPendingAccept, "user", order.UserID, "payment completed")
	return payment, cloneOrder(order), nil
}

func (s *MemoryStore) CancelOrder(orderID, userID uint64, reason string) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrUnauthorized
	}
	switch order.Status {
	case OrderStatusPendingPayment, OrderStatusPendingAccept, OrderStatusAccepted:
	default:
		return nil, ErrInvalidOrderStatus
	}
	order.CancelReason = reason
	s.transitionOrder(order, OrderStatusCanceled, "user", userID, reason)
	return cloneOrder(order), nil
}

func (s *MemoryStore) ProviderAcceptOrder(userID, orderID uint64) (*Order, error) {
	return s.transitionByProvider(userID, orderID, OrderStatusPendingAccept, OrderStatusAccepted, "provider accepted")
}

func (s *MemoryStore) ProviderDepartOrder(userID, orderID uint64) (*Order, error) {
	return s.transitionByProvider(userID, orderID, OrderStatusAccepted, OrderStatusDeparted, "provider departed")
}

func (s *MemoryStore) ProviderArriveOrder(userID, orderID uint64) (*Order, error) {
	return s.transitionByProvider(userID, orderID, OrderStatusDeparted, OrderStatusArrived, "provider arrived")
}

func (s *MemoryStore) StartOrder(userID, orderID uint64, remark string) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrUnauthorized
	}
	if order.Status != OrderStatusArrived {
		return nil, ErrInvalidOrderStatus
	}

	now := time.Now()
	order.StartedAt = &now
	s.transitionOrder(order, OrderStatusServing, "user", userID, remark)
	return cloneOrder(order), nil
}

func (s *MemoryStore) ProviderFinishOrder(userID, orderID uint64) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	provider, ok := s.providerByUID[userID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	if order.ProviderID != provider.ID {
		return nil, ErrUnauthorized
	}
	if order.Status != OrderStatusServing {
		return nil, ErrInvalidOrderStatus
	}

	now := time.Now()
	order.FinishedAt = &now
	s.transitionOrder(order, OrderStatusPendingFinish, "provider", userID, "provider finished service")
	return cloneOrder(order), nil
}

func (s *MemoryStore) ConfirmFinishOrder(userID, orderID uint64) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrUnauthorized
	}
	if order.Status != OrderStatusPendingFinish {
		return nil, ErrInvalidOrderStatus
	}
	s.transitionOrder(order, OrderStatusCompleted, "user", userID, "user confirmed finish")
	return cloneOrder(order), nil
}

func (s *MemoryStore) transitionByProvider(userID, orderID uint64, expectedStatus, targetStatus int, remark string) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	provider, ok := s.providerByUID[userID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	if provider.AuditStatus != 2 {
		return nil, ErrProviderAuditRequired
	}
	if order.ProviderID != provider.ID {
		return nil, ErrUnauthorized
	}
	if order.Status != expectedStatus {
		return nil, ErrInvalidOrderStatus
	}

	s.transitionOrder(order, targetStatus, "provider", userID, remark)
	return cloneOrder(order), nil
}

func (s *MemoryStore) transitionOrder(order *Order, toStatus int, role string, operatorID uint64, remark string) {
	now := time.Now()
	order.Logs = append(order.Logs, OrderLog{
		FromStatus:   order.Status,
		ToStatus:     toStatus,
		OperatorRole: role,
		OperatorID:   operatorID,
		Remark:       remark,
		CreatedAt:    now,
	})
	order.Status = toStatus
	order.UpdatedAt = now
}
