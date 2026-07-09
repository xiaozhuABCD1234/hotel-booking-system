// Package service 提供业务逻辑层，位于 handler 和 repo 之间。
package service

import (
	"context"
	"errors"
	"fmt"

	model "backend/model/schema"
	"backend/model/view"
	"backend/repo"

	"github.com/google/uuid"
)

// ErrInvalidTransition 非法状态流转错误。
var ErrInvalidTransition = errors.New("invalid status transition")

// OrderService 订单业务逻辑。
type OrderService struct {
	orders    repo.OrderRepository
	details   repo.OrderDetailRepository
	summaries repo.OrderSummaryRepository
}

// NewOrderService 创建 OrderService 实例。
func NewOrderService(
	orders repo.OrderRepository,
	details repo.OrderDetailRepository,
	summaries repo.OrderSummaryRepository,
) *OrderService {
	return &OrderService{orders: orders, details: details, summaries: summaries}
}

// Create 创建订单，校验基本字段。
func (s *OrderService) Create(ctx context.Context, order *model.Order) error {
	if order.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}
	if order.RoomID == uuid.Nil {
		return errors.New("room_id is required")
	}
	if order.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if !order.CheckOutDate.After(order.CheckInDate) {
		return errors.New("check_out_date must be after check_in_date")
	}
	if order.Status == "" {
		order.Status = model.OrderPending
	} else if order.Status != model.OrderPending {
		return errors.New("can only create orders with pending status")
	}

	return s.orders.Create(ctx, order)
}

// GetByID 查询订单详情。
func (s *OrderService) GetByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	return s.orders.FindByID(ctx, id)
}

// GetDetail 查询订单完整详情（走 fn_order_detail_1718，下单人/入住人明确区分）。
func (s *OrderService) GetDetail(ctx context.Context, id uuid.UUID) (*view.OrderDetail, error) {
	return s.details.FindDetailByOrderID(ctx, id)
}

// ListByUser 按用户查询订单列表。
func (s *OrderService) ListByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]model.Order, int64, error) {
	return s.orders.FindByUserID(ctx, userID, offset, limit)
}

// ListAll 查询全部订单。
func (s *OrderService) ListAll(ctx context.Context, offset, limit int) ([]model.Order, int64, error) {
	return s.orders.FindAll(ctx, offset, limit)
}

// ListByHotel 按酒店查询订单列表。
func (s *OrderService) ListByHotel(ctx context.Context, hotelID uuid.UUID, offset, limit int) ([]model.Order, int64, error) {
	return s.orders.FindByHotelID(ctx, hotelID, offset, limit)
}

// ListSummaries 查询订单概览列表，支持按状态筛选。
func (s *OrderService) ListSummaries(ctx context.Context, status string, offset, limit int) ([]view.OrderSummary, int64, error) {
	if status != "" {
		return s.summaries.FindByStatus(ctx, status, offset, limit)
	}
	return s.summaries.FindAll(ctx, offset, limit)
}

// UpdateStatus 更新订单状态，校验流转合法性。
//
// 合法流转：
//
//	pending  → booked, cancelled
//	booked   → checked_in, cancelled
//	checked_in → completed
//
// 终态 booked/completed/cancelled 不可再变更。
func (s *OrderService) UpdateStatus(ctx context.Context, id uuid.UUID, newStatus model.OrderStatus) error {
	order, err := s.orders.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// 终态不可变更
	switch order.Status {
	case model.OrderCompleted, model.OrderCancelled:
		return fmt.Errorf("%w: order is in terminal state %s", ErrInvalidTransition, order.Status)
	}

	// 状态机校验
	allowed := allowedTransitions(order.Status)
	for _, st := range allowed {
		if st == newStatus {
			return s.orders.UpdateStatus(ctx, id, newStatus)
		}
	}

	return fmt.Errorf("%w: cannot transition from %s to %s", ErrInvalidTransition, order.Status, newStatus)
}

// Delete 删除订单。
func (s *OrderService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.orders.Delete(ctx, id)
}

// allowedTransitions 返回当前状态允许流转到的目标状态列表。
func allowedTransitions(current model.OrderStatus) []model.OrderStatus {
	switch current {
	case model.OrderPending:
		return []model.OrderStatus{model.OrderBooked, model.OrderCancelled}
	case model.OrderBooked:
		return []model.OrderStatus{model.OrderCheckedIn, model.OrderCancelled}
	case model.OrderCheckedIn:
		return []model.OrderStatus{model.OrderCompleted}
	default:
		return nil
	}
}
