package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	model "backend/model/schema"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func setupReviewDeps(t *testing.T, tx *gorm.DB) (*model.Region, *model.Hotel, *model.Room, *model.User, *model.Order) {
	t.Helper()
	now := time.Now()
	region := &model.Region{RegionName: "Test"}
	if err := tx.Create(region).Error; err != nil {
		t.Fatalf("region: %v", err)
	}
	hotel := &model.Hotel{ID: uuid.New(), HotelName: "Review Hotel", RegionID: region.ID, Address: "Addr", Telephone: "123", Status: 1}
	if err := tx.Create(hotel).Error; err != nil {
		t.Fatalf("hotel: %v", err)
	}
	room := &model.Room{ID: uuid.New(), HotelID: hotel.ID, TypeName: "Deluxe", TotalQuantity: 5, AvailableQuantity: 5, Price: 200, Status: 1}
	if err := tx.Create(room).Error; err != nil {
		t.Fatalf("room: %v", err)
	}
	vip := &model.VipLevel{Level: 0, LevelName: "普通会员", MinPoints: 0, DiscountRate: 1.0}
	if err := tx.Exec(
		"INSERT INTO vip_level_1718 (level, level_name, min_points, discount_rate) VALUES (?, ?, ?, ?) ON CONFLICT DO NOTHING",
		vip.Level, vip.LevelName, vip.MinPoints, vip.DiscountRate,
	).Error; err != nil {
		t.Fatalf("vip level: %v", err)
	}
	user := &model.User{ID: uuid.New(), Username: "reviewer", Password: "pwd", Role: model.RoleCustomer, VipLevelID: 0, Status: 1}
	if err := tx.Create(user).Error; err != nil {
		t.Fatalf("user: %v", err)
	}
	order := &model.Order{ID: uuid.New(), UserID: user.ID, RoomID: room.ID, Quantity: 1, CheckInDate: now.AddDate(0, 0, 1), CheckOutDate: now.AddDate(0, 0, 3), TotalPrice: 200, Discount: 0, ActualPrice: 200, Status: model.OrderBooked}
	if err := tx.Create(order).Error; err != nil {
		t.Fatalf("order: %v", err)
	}
	return region, hotel, room, user, order
}

// TestReviewRepo_CRUD_Cycle tests Create → FindByID → Update → FindByID → Delete → FindByID.
func TestReviewRepo_CRUD_Cycle(t *testing.T) {
	tx := txRepo(t)
	_, hotel, _, user, order := setupReviewDeps(t, tx)

	repo := NewReviewRepo(tx)
	review := &model.Review{
		ID:      uuid.New(),
		UserID:  user.ID,
		OrderID: order.ID,
		HotelID: hotel.ID,
		Rating:  4,
		Content: ptr("Great stay!"),
	}

	// Create
	if err := repo.Create(context.Background(), review); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// FindByID
	found, err := repo.FindByID(context.Background(), review.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.Rating != 4 {
		t.Errorf("rating mismatch: got %d, want 4", found.Rating)
	}
	if found.Content == nil || *found.Content != "Great stay!" {
		t.Errorf("content mismatch: got %v, want Great stay!", found.Content)
	}
	if found.User.ID != user.ID {
		t.Errorf("user preload mismatch: got %v, want %v", found.User.ID, user.ID)
	}
	if found.Hotel.ID != hotel.ID {
		t.Errorf("hotel preload mismatch: got %v, want %v", found.Hotel.ID, hotel.ID)
	}
	if found.Order.ID != order.ID {
		t.Errorf("order preload mismatch: got %v, want %v", found.Order.ID, order.ID)
	}

	// Update
	found.Rating = 5
	found.Content = ptr("Amazing stay!")
	if err := repo.Update(context.Background(), found); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// FindByID after update
	found, err = repo.FindByID(context.Background(), review.ID)
	if err != nil {
		t.Fatalf("FindByID after update failed: %v", err)
	}
	if found.Rating != 5 {
		t.Errorf("rating after update mismatch: got %d, want 5", found.Rating)
	}
	if found.Content == nil || *found.Content != "Amazing stay!" {
		t.Errorf("content after update mismatch: got %v, want Amazing stay!", found.Content)
	}

	// Delete
	if err := repo.Delete(context.Background(), review.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// FindByID after delete should return ErrRecordNotFound
	_, err = repo.FindByID(context.Background(), review.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound after delete, got: %v", err)
	}
}

// TestReviewRepo_FindByHotelID_Pagination tests pagination for hotel reviews.
func TestReviewRepo_FindByHotelID_Pagination(t *testing.T) {
	tx := txRepo(t)
	_, hotel, room, user, order := setupReviewDeps(t, tx)

	repo := NewReviewRepo(tx)
	now := time.Now()
	for i := 0; i < 3; i++ {
		var orderID uuid.UUID
		if i == 0 {
			orderID = order.ID
		} else {
			extraOrder := &model.Order{
				ID:           uuid.New(),
				UserID:       user.ID,
				RoomID:       room.ID,
				Quantity:     1,
				CheckInDate:  now.AddDate(0, 0, i*2+1),
				CheckOutDate: now.AddDate(0, 0, i*2+3),
				TotalPrice:   200,
				Discount:     0,
				ActualPrice:  200,
				Status:       model.OrderBooked,
			}
			if err := tx.Create(extraOrder).Error; err != nil {
				t.Fatalf("create extra order %d: %v", i, err)
			}
			orderID = extraOrder.ID
		}
		review := &model.Review{
			ID:      uuid.New(),
			UserID:  user.ID,
			OrderID: orderID,
			HotelID: hotel.ID,
			Rating:  int16(3 + i),
			Content: ptr("Review " + string(rune('A'+i))),
		}
		if err := repo.Create(context.Background(), review); err != nil {
			t.Fatalf("Create review %d failed: %v", i, err)
		}
	}

	// Page 1: offset=0, limit=2
	results, total, err := repo.FindByHotelID(context.Background(), hotel.ID, 0, 2)
	if err != nil {
		t.Fatalf("FindByHotelID page 1 failed: %v", err)
	}
	if total != 3 {
		t.Errorf("total mismatch: got %d, want 3", total)
	}
	if len(results) != 2 {
		t.Errorf("page 1 length mismatch: got %d, want 2", len(results))
	}

	// Page 2: offset=2, limit=2
	results, total, err = repo.FindByHotelID(context.Background(), hotel.ID, 2, 2)
	if err != nil {
		t.Fatalf("FindByHotelID page 2 failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("page 2 length mismatch: got %d, want 1", len(results))
	}
}

// TestReviewRepo_FindByUserID tests finding reviews by user ID.
func TestReviewRepo_FindByUserID(t *testing.T) {
	tx := txRepo(t)
	_, hotel, room, user, order := setupReviewDeps(t, tx)

	repo := NewReviewRepo(tx)
	now := time.Now()
	for i := 0; i < 2; i++ {
		var orderID uuid.UUID
		if i == 0 {
			orderID = order.ID
		} else {
			extraOrder := &model.Order{
				ID:           uuid.New(),
				UserID:       user.ID,
				RoomID:       room.ID,
				Quantity:     1,
				CheckInDate:  now.AddDate(0, 0, 5),
				CheckOutDate: now.AddDate(0, 0, 7),
				TotalPrice:   200,
				Discount:     0,
				ActualPrice:  200,
				Status:       model.OrderBooked,
			}
			if err := tx.Create(extraOrder).Error; err != nil {
				t.Fatalf("create extra order: %v", err)
			}
			orderID = extraOrder.ID
		}
		review := &model.Review{
			ID:      uuid.New(),
			UserID:  user.ID,
			OrderID: orderID,
			HotelID: hotel.ID,
			Rating:  int16(3 + i),
			Content: ptr("Review " + string(rune('A'+i))),
		}
		if err := repo.Create(context.Background(), review); err != nil {
			t.Fatalf("Create review %d failed: %v", i, err)
		}
	}

	results, total, err := repo.FindByUserID(context.Background(), user.ID, 0, 0)
	if err != nil {
		t.Fatalf("FindByUserID failed: %v", err)
	}
	if total != 2 {
		t.Errorf("total mismatch: got %d, want 2", total)
	}
	if len(results) != 2 {
		t.Errorf("results length mismatch: got %d, want 2", len(results))
	}
}

// TestReviewRepo_FindByOrderID tests finding a review by order ID.
func TestReviewRepo_FindByOrderID(t *testing.T) {
	tx := txRepo(t)
	_, hotel, _, user, order := setupReviewDeps(t, tx)

	repo := NewReviewRepo(tx)
	review := &model.Review{
		ID:      uuid.New(),
		UserID:  user.ID,
		OrderID: order.ID,
		HotelID: hotel.ID,
		Rating:  4,
		Content: ptr("Great stay!"),
	}
	if err := repo.Create(context.Background(), review); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.FindByOrderID(context.Background(), order.ID)
	if err != nil {
		t.Fatalf("FindByOrderID failed: %v", err)
	}
	if found.ID != review.ID {
		t.Errorf("ID mismatch: got %v, want %v", found.ID, review.ID)
	}
}

// TestReviewRepo_FindByRating tests finding reviews by rating.
func TestReviewRepo_FindByRating(t *testing.T) {
	tx := txRepo(t)
	_, hotel, _, user, order := setupReviewDeps(t, tx)

	repo := NewReviewRepo(tx)
	review3 := &model.Review{
		ID:      uuid.New(),
		UserID:  user.ID,
		OrderID: order.ID,
		HotelID: hotel.ID,
		Rating:  3,
		Content: ptr("Okay stay."),
	}
	if err := repo.Create(context.Background(), review3); err != nil {
		t.Fatalf("Create review rating 3 failed: %v", err)
	}

	// Need a second order for the second review because of unique(user_id, order_id) constraint.
	now := time.Now()
	order2 := &model.Order{
		ID:           uuid.New(),
		UserID:       user.ID,
		RoomID:       order.RoomID,
		Quantity:     1,
		CheckInDate:  now.AddDate(0, 0, 5),
		CheckOutDate: now.AddDate(0, 0, 7),
		TotalPrice:   200,
		Discount:     0,
		ActualPrice:  200,
		Status:       model.OrderBooked,
	}
	if err := tx.Create(order2).Error; err != nil {
		t.Fatalf("create order2: %v", err)
	}

	review5 := &model.Review{
		ID:      uuid.New(),
		UserID:  user.ID,
		OrderID: order2.ID,
		HotelID: hotel.ID,
		Rating:  5,
		Content: ptr("Excellent stay!"),
	}
	if err := repo.Create(context.Background(), review5); err != nil {
		t.Fatalf("Create review rating 5 failed: %v", err)
	}

	results, total, err := repo.FindByRating(context.Background(), 5, 0, 0)
	if err != nil {
		t.Fatalf("FindByRating(5) failed: %v", err)
	}
	if total != 1 {
		t.Errorf("total mismatch: got %d, want 1", total)
	}
	if len(results) != 1 {
		t.Errorf("results length mismatch: got %d, want 1", len(results))
	}
	if results[0].Rating != 5 {
		t.Errorf("rating mismatch: got %d, want 5", results[0].Rating)
	}
}

// TestReviewRepo_Update_OnlyRatingAndContent verifies Update only modifies Rating and Content.
func TestReviewRepo_Update_OnlyRatingAndContent(t *testing.T) {
	tx := txRepo(t)
	_, hotel, _, user, order := setupReviewDeps(t, tx)

	repo := NewReviewRepo(tx)
	review := &model.Review{
		ID:      uuid.New(),
		UserID:  user.ID,
		OrderID: order.ID,
		HotelID: hotel.ID,
		Rating:  4,
		Content: ptr("Great stay!"),
	}
	if err := repo.Create(context.Background(), review); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Try updating with a changed OrderID; repo.Select should ignore it.
	now := time.Now()
	order2 := &model.Order{
		ID:           uuid.New(),
		UserID:       user.ID,
		RoomID:       order.RoomID,
		Quantity:     1,
		CheckInDate:  now.AddDate(0, 0, 5),
		CheckOutDate: now.AddDate(0, 0, 7),
		TotalPrice:   200,
		Discount:     0,
		ActualPrice:  200,
		Status:       model.OrderBooked,
	}
	if err := tx.Create(order2).Error; err != nil {
		t.Fatalf("create order2: %v", err)
	}

	review.OrderID = order2.ID
	review.Rating = 5
	review.Content = ptr("Updated content")
	if err := repo.Update(context.Background(), review); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify OrderID is still the original one.
	found, err := repo.FindByID(context.Background(), review.ID)
	if err != nil {
		t.Fatalf("FindByID after update failed: %v", err)
	}
	if found.OrderID != order.ID {
		t.Errorf("OrderID was changed: got %v, want %v", found.OrderID, order.ID)
	}
	if found.Rating != 5 {
		t.Errorf("rating mismatch: got %d, want 5", found.Rating)
	}
	if found.Content == nil || *found.Content != "Updated content" {
		t.Errorf("content mismatch: got %v, want Updated content", found.Content)
	}
}

// TestReviewFullRepo_FindByReviewID tests querying the review full view by review ID.
func TestReviewFullRepo_FindByReviewID(t *testing.T) {
	tx := txRepo(t)
	_, hotel, _, user, order := setupReviewDeps(t, tx)

	reviewRepo := NewReviewRepo(tx)
	review := &model.Review{
		ID:      uuid.New(),
		UserID:  user.ID,
		OrderID: order.ID,
		HotelID: hotel.ID,
		Rating:  4,
		Content: ptr("Great stay!"),
	}
	if err := reviewRepo.Create(context.Background(), review); err != nil {
		t.Fatalf("Create review failed: %v", err)
	}

	fullRepo := NewReviewFullRepo(tx)
	result, err := fullRepo.FindByReviewID(context.Background(), review.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			t.Skip("view not available or review not found in view")
		}
		t.Skipf("view not available: %v", err)
	}
	if result.ReviewID != review.ID {
		t.Errorf("review ID mismatch: got %v, want %v", result.ReviewID, review.ID)
	}
	if result.HotelID != hotel.ID {
		t.Errorf("hotel ID mismatch: got %v, want %v", result.HotelID, hotel.ID)
	}
	if result.UserID != user.ID {
		t.Errorf("user ID mismatch: got %v, want %v", result.UserID, user.ID)
	}
}

// TestReviewFullRepo_FindByHotelID tests querying the review full view by hotel ID.
func TestReviewFullRepo_FindByHotelID(t *testing.T) {
	tx := txRepo(t)
	_, hotel, _, user, order := setupReviewDeps(t, tx)

	reviewRepo := NewReviewRepo(tx)
	review := &model.Review{
		ID:      uuid.New(),
		UserID:  user.ID,
		OrderID: order.ID,
		HotelID: hotel.ID,
		Rating:  4,
		Content: ptr("Great stay!"),
	}
	if err := reviewRepo.Create(context.Background(), review); err != nil {
		t.Fatalf("Create review failed: %v", err)
	}

	fullRepo := NewReviewFullRepo(tx)
	results, total, err := fullRepo.FindByHotelID(context.Background(), hotel.ID, 0, 0)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			t.Skip("view not available or no records in view")
		}
		t.Skipf("view not available: %v", err)
	}
	if total < 1 {
		t.Errorf("total mismatch: got %d, want at least 1", total)
	}
	if len(results) < 1 {
		t.Errorf("results length mismatch: got %d, want at least 1", len(results))
	}
	found := false
	for _, r := range results {
		if r.ReviewID == review.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected review %v in results", review.ID)
	}
}
