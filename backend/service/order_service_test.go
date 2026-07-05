package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	model "backend/model/schema"
	"backend/repo"
	"backend/service"

	"github.com/google/uuid"
)

// ─── Mock Repository ───────────────────────────────────────────

type mockOrderRepo struct {
	createFunc                func(ctx context.Context, order *model.Order) error
	findByIDFunc              func(ctx context.Context, id uuid.UUID) (*model.Order, error)
	findByUserIDFunc          func(ctx context.Context, userID uuid.UUID, offset, limit int) ([]model.Order, int64, error)
	findByUserIDAndStatusFunc func(ctx context.Context, userID uuid.UUID, status model.OrderStatus, offset, limit int) ([]model.Order, int64, error)
	findByHotelIDFunc         func(ctx context.Context, hotelID uuid.UUID, offset, limit int) ([]model.Order, int64, error)
	findAllFunc               func(ctx context.Context, offset, limit int) ([]model.Order, int64, error)
	updateFunc                func(ctx context.Context, order *model.Order) error
	updateStatusFunc          func(ctx context.Context, id uuid.UUID, status model.OrderStatus) error
	deleteFunc                func(ctx context.Context, id uuid.UUID) error
}

var _ repo.OrderRepository = (*mockOrderRepo)(nil)

func (m *mockOrderRepo) Create(ctx context.Context, order *model.Order) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, order)
	}
	return nil
}

func (m *mockOrderRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockOrderRepo) FindByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]model.Order, int64, error) {
	if m.findByUserIDFunc != nil {
		return m.findByUserIDFunc(ctx, userID, offset, limit)
	}
	return nil, 0, nil
}

func (m *mockOrderRepo) FindByUserIDAndStatus(ctx context.Context, userID uuid.UUID, status model.OrderStatus, offset, limit int) ([]model.Order, int64, error) {
	if m.findByUserIDAndStatusFunc != nil {
		return m.findByUserIDAndStatusFunc(ctx, userID, status, offset, limit)
	}
	return nil, 0, nil
}

func (m *mockOrderRepo) FindByHotelID(ctx context.Context, hotelID uuid.UUID, offset, limit int) ([]model.Order, int64, error) {
	if m.findByHotelIDFunc != nil {
		return m.findByHotelIDFunc(ctx, hotelID, offset, limit)
	}
	return nil, 0, nil
}

func (m *mockOrderRepo) FindAll(ctx context.Context, offset, limit int) ([]model.Order, int64, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx, offset, limit)
	}
	return nil, 0, nil
}

func (m *mockOrderRepo) Update(ctx context.Context, order *model.Order) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, order)
	}
	return nil
}

func (m *mockOrderRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) error {
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(ctx, id, status)
	}
	return nil
}

func (m *mockOrderRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

// ─── Helpers ─────────────────────────────────────────────────

func newValidOrder() *model.Order {
	return &model.Order{
		ID:           uuid.New(),
		UserID:       uuid.New(),
		RoomID:       uuid.New(),
		Quantity:     1,
		CheckInDate:  time.Now().AddDate(0, 0, 1),
		CheckOutDate: time.Now().AddDate(0, 0, 3),
		TotalPrice:   200,
		Discount:     0,
		ActualPrice:  200,
		Status:       model.OrderPending,
		Guests:       []model.OrderGuest{},
	}
}

func orderWithStatus(status model.OrderStatus) *model.Order {
	order := newValidOrder()
	order.Status = status
	return order
}

// ─── Tests: Create ─────────────────────────────────────────────

func TestOrderService_Create(t *testing.T) {
	t.Run("nil userID", func(t *testing.T) {
		mock := &mockOrderRepo{}
		svc := service.NewOrderService(mock)
		order := newValidOrder()
		order.UserID = uuid.Nil

		err := svc.Create(context.Background(), order)
		if err == nil {
			t.Fatalf("expected error for nil userID, got nil")
		}
		if err.Error() != "user_id is required" {
			t.Errorf("unexpected error message: %q", err.Error())
		}
	})

	t.Run("nil roomID", func(t *testing.T) {
		mock := &mockOrderRepo{}
		svc := service.NewOrderService(mock)
		order := newValidOrder()
		order.RoomID = uuid.Nil

		err := svc.Create(context.Background(), order)
		if err == nil {
			t.Fatalf("expected error for nil roomID, got nil")
		}
		if err.Error() != "room_id is required" {
			t.Errorf("unexpected error message: %q", err.Error())
		}
	})

	t.Run("zero quantity", func(t *testing.T) {
		mock := &mockOrderRepo{}
		svc := service.NewOrderService(mock)
		order := newValidOrder()
		order.Quantity = 0

		err := svc.Create(context.Background(), order)
		if err == nil {
			t.Fatalf("expected error for zero quantity, got nil")
		}
		if err.Error() != "quantity must be positive" {
			t.Errorf("unexpected error message: %q", err.Error())
		}
	})

	t.Run("negative quantity", func(t *testing.T) {
		mock := &mockOrderRepo{}
		svc := service.NewOrderService(mock)
		order := newValidOrder()
		order.Quantity = -1

		err := svc.Create(context.Background(), order)
		if err == nil {
			t.Fatalf("expected error for negative quantity, got nil")
		}
		if err.Error() != "quantity must be positive" {
			t.Errorf("unexpected error message: %q", err.Error())
		}
	})

	t.Run("invalid dates", func(t *testing.T) {
		mock := &mockOrderRepo{}
		svc := service.NewOrderService(mock)
		order := newValidOrder()
		order.CheckOutDate = order.CheckInDate

		err := svc.Create(context.Background(), order)
		if err == nil {
			t.Fatalf("expected error for invalid dates, got nil")
		}
		if err.Error() != "check_out_date must be after check_in_date" {
			t.Errorf("unexpected error message: %q", err.Error())
		}
	})

	t.Run("check_out before check_in", func(t *testing.T) {
		mock := &mockOrderRepo{}
		svc := service.NewOrderService(mock)
		order := newValidOrder()
		order.CheckOutDate = order.CheckInDate.AddDate(0, 0, -1)

		err := svc.Create(context.Background(), order)
		if err == nil {
			t.Fatalf("expected error for check_out before check_in, got nil")
		}
		if err.Error() != "check_out_date must be after check_in_date" {
			t.Errorf("unexpected error message: %q", err.Error())
		}
	})

	t.Run("success", func(t *testing.T) {
		created := false
		mock := &mockOrderRepo{
			createFunc: func(ctx context.Context, order *model.Order) error {
				created = true
				return nil
			},
		}
		svc := service.NewOrderService(mock)
		order := newValidOrder()

		err := svc.Create(context.Background(), order)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !created {
			t.Errorf("expected repo.Create to be called")
		}
	})

	t.Run("repo error propagated", func(t *testing.T) {
		repoErr := errors.New("db failure")
		mock := &mockOrderRepo{
			createFunc: func(ctx context.Context, order *model.Order) error {
				return repoErr
			},
		}
		svc := service.NewOrderService(mock)
		order := newValidOrder()

		err := svc.Create(context.Background(), order)
		if !errors.Is(err, repoErr) {
			t.Errorf("expected repo error to be propagated, got: %v", err)
		}
	})
}

// ─── Tests: UpdateStatus ───────────────────────────────────────

func TestOrderService_UpdateStatus(t *testing.T) {
	makeSvc := func(status model.OrderStatus) (*service.OrderService, *mockOrderRepo) {
		mock := &mockOrderRepo{
			findByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Order, error) {
				return orderWithStatus(status), nil
			},
			updateStatusFunc: func(ctx context.Context, id uuid.UUID, s model.OrderStatus) error {
				return nil
			},
		}
		return service.NewOrderService(mock), mock
	}

	// Legal transitions
	t.Run("pending to booked", func(t *testing.T) {
		svc, _ := makeSvc(model.OrderPending)
		err := svc.UpdateStatus(context.Background(), uuid.New(), model.OrderBooked)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("pending to cancelled", func(t *testing.T) {
		svc, _ := makeSvc(model.OrderPending)
		err := svc.UpdateStatus(context.Background(), uuid.New(), model.OrderCancelled)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("booked to checked_in", func(t *testing.T) {
		svc, _ := makeSvc(model.OrderBooked)
		err := svc.UpdateStatus(context.Background(), uuid.New(), model.OrderCheckedIn)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("booked to cancelled", func(t *testing.T) {
		svc, _ := makeSvc(model.OrderBooked)
		err := svc.UpdateStatus(context.Background(), uuid.New(), model.OrderCancelled)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("checked_in to completed", func(t *testing.T) {
		svc, _ := makeSvc(model.OrderCheckedIn)
		err := svc.UpdateStatus(context.Background(), uuid.New(), model.OrderCompleted)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	// Illegal transitions
	t.Run("booked to pending", func(t *testing.T) {
		svc, _ := makeSvc(model.OrderBooked)
		err := svc.UpdateStatus(context.Background(), uuid.New(), model.OrderPending)
		if err == nil {
			t.Fatalf("expected error for booked→pending, got nil")
		}
		if !errors.Is(err, service.ErrInvalidTransition) {
			t.Errorf("expected ErrInvalidTransition, got: %v", err)
		}
	})

	t.Run("checked_in to pending", func(t *testing.T) {
		svc, _ := makeSvc(model.OrderCheckedIn)
		err := svc.UpdateStatus(context.Background(), uuid.New(), model.OrderPending)
		if err == nil {
			t.Fatalf("expected error for checked_in→pending, got nil")
		}
		if !errors.Is(err, service.ErrInvalidTransition) {
			t.Errorf("expected ErrInvalidTransition, got: %v", err)
		}
	})

	t.Run("checked_in to booked", func(t *testing.T) {
		svc, _ := makeSvc(model.OrderCheckedIn)
		err := svc.UpdateStatus(context.Background(), uuid.New(), model.OrderBooked)
		if err == nil {
			t.Fatalf("expected error for checked_in→booked, got nil")
		}
		if !errors.Is(err, service.ErrInvalidTransition) {
			t.Errorf("expected ErrInvalidTransition, got: %v", err)
		}
	})

	t.Run("checked_in to cancelled", func(t *testing.T) {
		svc, _ := makeSvc(model.OrderCheckedIn)
		err := svc.UpdateStatus(context.Background(), uuid.New(), model.OrderCancelled)
		if err == nil {
			t.Fatalf("expected error for checked_in→cancelled, got nil")
		}
		if !errors.Is(err, service.ErrInvalidTransition) {
			t.Errorf("expected ErrInvalidTransition, got: %v", err)
		}
	})

	t.Run("completed to anything", func(t *testing.T) {
		targets := []model.OrderStatus{model.OrderPending, model.OrderBooked, model.OrderCheckedIn, model.OrderCancelled}
		for _, target := range targets {
			mock := &mockOrderRepo{
				findByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Order, error) {
					return orderWithStatus(model.OrderCompleted), nil
				},
			}
			svc := service.NewOrderService(mock)
			err := svc.UpdateStatus(context.Background(), uuid.New(), target)
			if err == nil {
				t.Fatalf("expected error for completed→%s, got nil", target)
			}
			if !errors.Is(err, service.ErrInvalidTransition) {
				t.Errorf("expected ErrInvalidTransition for completed→%s, got: %v", target, err)
			}
		}
	})

	t.Run("cancelled to anything", func(t *testing.T) {
		targets := []model.OrderStatus{model.OrderPending, model.OrderBooked, model.OrderCheckedIn, model.OrderCompleted}
		for _, target := range targets {
			mock := &mockOrderRepo{
				findByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Order, error) {
					return orderWithStatus(model.OrderCancelled), nil
				},
			}
			svc := service.NewOrderService(mock)
			err := svc.UpdateStatus(context.Background(), uuid.New(), target)
			if err == nil {
				t.Fatalf("expected error for cancelled→%s, got nil", target)
			}
			if !errors.Is(err, service.ErrInvalidTransition) {
				t.Errorf("expected ErrInvalidTransition for cancelled→%s, got: %v", target, err)
			}
		}
	})

	t.Run("order not found", func(t *testing.T) {
		repoErr := errors.New("record not found")
		mock := &mockOrderRepo{
			findByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Order, error) {
				return nil, repoErr
			},
		}
		svc := service.NewOrderService(mock)
		err := svc.UpdateStatus(context.Background(), uuid.New(), model.OrderBooked)
		if err == nil {
			t.Fatalf("expected error when order not found, got nil")
		}
		if errors.Is(err, service.ErrInvalidTransition) {
			t.Errorf("expected non-ErrInvalidTransition error, got: %v", err)
		}
	})
}

// ─── Tests: GetByID ────────────────────────────────────────────

func TestOrderService_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := newValidOrder()
		mock := &mockOrderRepo{
			findByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Order, error) {
				return expected, nil
			},
		}
		svc := service.NewOrderService(mock)
		order, err := svc.GetByID(context.Background(), expected.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if order.ID != expected.ID {
			t.Errorf("ID mismatch: got %v, want %v", order.ID, expected.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repoErr := errors.New("not found")
		mock := &mockOrderRepo{
			findByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Order, error) {
				return nil, repoErr
			},
		}
		svc := service.NewOrderService(mock)
		_, err := svc.GetByID(context.Background(), uuid.New())
		if !errors.Is(err, repoErr) {
			t.Errorf("expected repo error to be propagated, got: %v", err)
		}
	})
}

// ─── Tests: ListByUser ─────────────────────────────────────────

func TestOrderService_ListByUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		expected := []model.Order{*newValidOrder(), *newValidOrder()}
		mock := &mockOrderRepo{
			findByUserIDFunc: func(ctx context.Context, uid uuid.UUID, offset, limit int) ([]model.Order, int64, error) {
				return expected, int64(len(expected)), nil
			},
		}
		svc := service.NewOrderService(mock)
		orders, total, err := svc.ListByUser(context.Background(), userID, 0, 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if total != 2 {
			t.Errorf("total mismatch: got %d, want 2", total)
		}
		if len(orders) != 2 {
			t.Errorf("results length mismatch: got %d, want 2", len(orders))
		}
	})

	t.Run("repo error propagated", func(t *testing.T) {
		repoErr := errors.New("db failure")
		mock := &mockOrderRepo{
			findByUserIDFunc: func(ctx context.Context, uid uuid.UUID, offset, limit int) ([]model.Order, int64, error) {
				return nil, 0, repoErr
			},
		}
		svc := service.NewOrderService(mock)
		_, _, err := svc.ListByUser(context.Background(), uuid.New(), 0, 10)
		if !errors.Is(err, repoErr) {
			t.Errorf("expected repo error to be propagated, got: %v", err)
		}
	})
}

// ─── Tests: ListAll ────────────────────────────────────────────

func TestOrderService_ListAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []model.Order{*newValidOrder(), *newValidOrder(), *newValidOrder()}
		mock := &mockOrderRepo{
			findAllFunc: func(ctx context.Context, offset, limit int) ([]model.Order, int64, error) {
				return expected, int64(len(expected)), nil
			},
		}
		svc := service.NewOrderService(mock)
		orders, total, err := svc.ListAll(context.Background(), 0, 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if total != 3 {
			t.Errorf("total mismatch: got %d, want 3", total)
		}
		if len(orders) != 3 {
			t.Errorf("results length mismatch: got %d, want 3", len(orders))
		}
	})

	t.Run("repo error propagated", func(t *testing.T) {
		repoErr := errors.New("db failure")
		mock := &mockOrderRepo{
			findAllFunc: func(ctx context.Context, offset, limit int) ([]model.Order, int64, error) {
				return nil, 0, repoErr
			},
		}
		svc := service.NewOrderService(mock)
		_, _, err := svc.ListAll(context.Background(), 0, 10)
		if !errors.Is(err, repoErr) {
			t.Errorf("expected repo error to be propagated, got: %v", err)
		}
	})
}

// ─── Tests: Delete ─────────────────────────────────────────────

func TestOrderService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		deleted := false
		mock := &mockOrderRepo{
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				deleted = true
				return nil
			},
		}
		svc := service.NewOrderService(mock)
		err := svc.Delete(context.Background(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !deleted {
			t.Errorf("expected repo.Delete to be called")
		}
	})

	t.Run("repo error propagated", func(t *testing.T) {
		repoErr := errors.New("db failure")
		mock := &mockOrderRepo{
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return repoErr
			},
		}
		svc := service.NewOrderService(mock)
		err := svc.Delete(context.Background(), uuid.New())
		if !errors.Is(err, repoErr) {
			t.Errorf("expected repo error to be propagated, got: %v", err)
		}
	})
}
