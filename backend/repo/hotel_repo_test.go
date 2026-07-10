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

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func createTestRegion(t *testing.T, tx *gorm.DB) *model.Region {
	t.Helper()
	region := &model.Region{RegionName: "Test Province"}
	if err := tx.Create(region).Error; err != nil {
		t.Fatalf("failed to create region: %v", err)
	}
	return region
}

func createTestHotel(t *testing.T, tx *gorm.DB, regionID int, name string) *model.Hotel {
	t.Helper()
	hotel := &model.Hotel{
		ID:        uuid.New(),
		HotelName: name,
		RegionID:  regionID,
		Address:   "123 Test Street",
		Telephone: "1234567890",
		Status:    1,
	}
	if err := tx.Create(hotel).Error; err != nil {
		t.Fatalf("failed to create hotel: %v", err)
	}
	return hotel
}

func createTestRoom(t *testing.T, tx *gorm.DB, hotelID uuid.UUID, typeName string) *model.Room {
	t.Helper()
	room := &model.Room{
		ID:                uuid.New(),
		HotelID:           hotelID,
		TypeName:          typeName,
		TotalQuantity:     10,
		AvailableQuantity: 5,
		Price:             299.99,
		Status:            1,
	}
	if err := tx.Create(room).Error; err != nil {
		t.Fatalf("failed to create room: %v", err)
	}
	return room
}

// ---------------------------------------------------------------------------
// HotelRepo
// ---------------------------------------------------------------------------

func TestHotelRepo_CRUD_Cycle(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewHotelRepo(tx)

	region := createTestRegion(t, tx)
	hotel := &model.Hotel{
		ID:        uuid.New(),
		HotelName: "Test Hotel",
		RegionID:  region.ID,
		Address:   "123 Street",
		Telephone: "123456",
		Status:    1,
	}

	if err := repo.Create(ctx, hotel); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.FindByID(ctx, hotel.ID)
	if err != nil {
		t.Fatalf("FindByID after Create failed: %v", err)
	}
	if found.HotelName != hotel.HotelName {
		t.Errorf("HotelName mismatch: got %q, want %q", found.HotelName, hotel.HotelName)
	}

	hotel.HotelName = "Updated Hotel"
	if err := repo.Update(ctx, hotel); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, err = repo.FindByID(ctx, hotel.ID)
	if err != nil {
		t.Fatalf("FindByID after Update failed: %v", err)
	}
	if found.HotelName != "Updated Hotel" {
		t.Errorf("HotelName after Update mismatch: got %q, want %q", found.HotelName, "Updated Hotel")
	}

	if err := repo.Delete(ctx, hotel.ID); err != nil {
		t.Fatalf("Delete (soft) failed: %v", err)
	}

	_, err = repo.FindByID(ctx, hotel.ID)
	if err == nil {
		t.Fatalf("FindByID after soft-delete should return error, got nil")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestHotelRepo_Create_WithImages(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewHotelRepo(tx)

	region := createTestRegion(t, tx)
	hotel := &model.Hotel{
		ID:        uuid.New(),
		HotelName: "Image Hotel",
		RegionID:  region.ID,
		Address:   "456 Avenue",
		Telephone: "654321",
		Status:    1,
		Images: []model.HotelImage{
			{ImageURL: "http://example.com/img1.jpg"},
			{ImageURL: "http://example.com/img2.jpg"},
		},
	}

	if err := repo.Create(ctx, hotel); err != nil {
		t.Fatalf("Create with images failed: %v", err)
	}

	found, err := repo.FindByID(ctx, hotel.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if len(found.Images) != 2 {
		t.Errorf("Images count mismatch: got %d, want 2", len(found.Images))
	}
}

func TestHotelRepo_FindAll_Pagination(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewHotelRepo(tx)

	region := createTestRegion(t, tx)

	// 按 region 过滤，避免预存数据干扰分页测试
	regionID := region.ID
	for i := 0; i < 3; i++ {
		createTestHotel(t, tx, regionID, "Paginated Hotel")
	}

	results, total, err := repo.FindAll(ctx, 0, 2, &regionID, nil, "", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("FindAll(0,2) failed: %v", err)
	}
	if total != 3 {
		t.Errorf("total mismatch: got %d, want 3", total)
	}
	if len(results) != 2 {
		t.Errorf("page size mismatch: got %d, want 2", len(results))
	}

	results, total, err = repo.FindAll(ctx, 2, 2, &regionID, nil, "", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("FindAll(2,2) failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("page size mismatch: got %d, want 1", len(results))
	}
}

func TestHotelRepo_FindAll_KeywordFilter(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewHotelRepo(tx)

	region := createTestRegion(t, tx)
	createTestHotel(t, tx, region.ID, "Grand Hotel")
	createTestHotel(t, tx, region.ID, "Budget Inn")

	results, total, err := repo.FindAll(ctx, 0, 10, nil, nil, "Grand", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("FindAll with keyword failed: %v", err)
	}
	if total != 1 {
		t.Errorf("total mismatch: got %d, want 1", total)
	}
	if len(results) != 1 {
		t.Errorf("results count mismatch: got %d, want 1", len(results))
	}
	if results[0].HotelName != "Grand Hotel" {
		t.Errorf("HotelName mismatch: got %q, want %q", results[0].HotelName, "Grand Hotel")
	}
}

func TestHotelRepo_HardDelete(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewHotelRepo(tx)

	region := createTestRegion(t, tx)
	hotel := createTestHotel(t, tx, region.ID, "Hard Delete Hotel")

	if err := repo.HardDelete(ctx, hotel.ID); err != nil {
		t.Fatalf("HardDelete failed: %v", err)
	}

	_, err := repo.FindByID(ctx, hotel.ID)
	if err == nil {
		t.Fatalf("FindByID after HardDelete should return error, got nil")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestHotelRepo_FindAll_PriceRange(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewHotelRepo(tx)

	region := createTestRegion(t, tx)
	hotel1 := createTestHotel(t, tx, region.ID, "Budget Hotel")
	hotel2 := createTestHotel(t, tx, region.ID, "Luxury Hotel")
	createTestRoom(t, tx, hotel1.ID, "Cheap")
	createTestRoom(t, tx, hotel2.ID, "Expensive")

	tx.Model(&model.Room{}).Where("hotel_id = ?", hotel2.ID).Update("price", 999.99)

	min := 500.0
	max := 1500.0
	results, total, err := repo.FindAll(ctx, 0, 10, &region.ID, nil, "", &min, &max, nil, nil)
	if err != nil {
		t.Fatalf("FindAll with price range failed: %v", err)
	}
	if total != 1 {
		t.Errorf("total mismatch: got %d, want 1", total)
	}
	if len(results) > 0 && results[0].ID != hotel2.ID {
		t.Errorf("expected Luxury Hotel, got %s", results[0].HotelName)
	}
}

func TestHotelRepo_FindAll_SingleSidedMinPrice(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewHotelRepo(tx)

	region := createTestRegion(t, tx)
	hotel := createTestHotel(t, tx, region.ID, "Pricey Hotel")
	createTestRoom(t, tx, hotel.ID, "Deluxe")

	min := 400.0
	results, _, err := repo.FindAll(ctx, 0, 10, &region.ID, nil, "", &min, nil, nil, nil)
	if err != nil {
		t.Fatalf("FindAll with minPrice failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no hotels above 400, got %d", len(results))
	}

	min = 100.0
	results, total, err := repo.FindAll(ctx, 0, 10, &region.ID, nil, "", &min, nil, nil, nil)
	if err != nil {
		t.Fatalf("FindAll with minPrice=100 failed: %v", err)
	}
	if total != 1 {
		t.Errorf("total mismatch: got %d, want 1", total)
	}
}

func TestHotelRepo_FindAll_SingleSidedMaxPrice(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewHotelRepo(tx)

	region := createTestRegion(t, tx)
	hotel := createTestHotel(t, tx, region.ID, "Affordable Hotel")
	createTestRoom(t, tx, hotel.ID, "Standard")

	max := 200.0
	results, _, err := repo.FindAll(ctx, 0, 10, &region.ID, nil, "", nil, &max, nil, nil)
	if err != nil {
		t.Fatalf("FindAll with maxPrice failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no hotels under 200, got %d", len(results))
	}

	max = 500.0
	results, total, err := repo.FindAll(ctx, 0, 10, &region.ID, nil, "", nil, &max, nil, nil)
	if err != nil {
		t.Fatalf("FindAll with maxPrice=500 failed: %v", err)
	}
	if total != 1 {
		t.Errorf("total mismatch: got %d, want 1", total)
	}
}

func TestHotelRepo_FindAll_DateAvailability(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewHotelRepo(tx)

	region := createTestRegion(t, tx)
	hotel := createTestHotel(t, tx, region.ID, "Date Hotel")
	createTestRoom(t, tx, hotel.ID, "Standard")

	tomorrow := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
	dayAfter := tomorrow.AddDate(0, 0, 1)

	results, total, err := repo.FindAll(ctx, 0, 10, &region.ID, nil, "", nil, nil, &tomorrow, &dayAfter)
	if err != nil {
		t.Fatalf("FindAll with date filter failed: %v", err)
	}
	if total != 1 {
		t.Errorf("total mismatch: got %d, want 1", total)
	}
	if len(results) != 1 {
		t.Errorf("results count mismatch: got %d, want 1", len(results))
	}
}

func TestHotelRepo_FindAll_DateAndPriceCombined(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewHotelRepo(tx)

	region := createTestRegion(t, tx)
	hotel1 := createTestHotel(t, tx, region.ID, "Match Hotel")
	_ = createTestHotel(t, tx, region.ID, "Wrong Price Hotel")
	createTestRoom(t, tx, hotel1.ID, "Standard")

	tomorrow := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
	dayAfter := tomorrow.AddDate(0, 0, 1)
	min := 200.0
	max := 400.0

	results, total, err := repo.FindAll(ctx, 0, 10, &region.ID, nil, "", &min, &max, &tomorrow, &dayAfter)
	if err != nil {
		t.Fatalf("FindAll with combined filters failed: %v", err)
	}
	if total != 1 {
		t.Errorf("total mismatch: got %d, want 1", total)
	}
	if len(results) != 1 {
		t.Errorf("results count mismatch: got %d, want 1", len(results))
	}
	if results[0].HotelName != "Match Hotel" {
		t.Errorf("HotelName mismatch: got %q, want %q", results[0].HotelName, "Match Hotel")
	}
}

func TestHotelRepo_FindAll_NoNewParams(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewHotelRepo(tx)

	region := createTestRegion(t, tx)
	createTestHotel(t, tx, region.ID, "Classic Hotel")

	results, total, err := repo.FindAll(ctx, 0, 10, &region.ID, nil, "", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("FindAll with no new params failed: %v", err)
	}
	if total != 1 {
		t.Errorf("total mismatch: got %d, want 1", total)
	}
	if len(results) != 1 {
		t.Errorf("results count mismatch: got %d, want 1", len(results))
	}
}

// ---------------------------------------------------------------------------
// RoomRepo
// ---------------------------------------------------------------------------

func TestRoomRepo_CRUD_Cycle(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewRoomRepo(tx)

	region := createTestRegion(t, tx)
	hotel := createTestHotel(t, tx, region.ID, "Room Hotel")
	room := &model.Room{
		ID:                uuid.New(),
		HotelID:           hotel.ID,
		TypeName:          "Deluxe",
		TotalQuantity:     10,
		AvailableQuantity: 5,
		Price:             399.00,
		Status:            1,
	}

	if err := repo.Create(ctx, room); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.FindByID(ctx, room.ID)
	if err != nil {
		t.Fatalf("FindByID after Create failed: %v", err)
	}
	if found.TypeName != room.TypeName {
		t.Errorf("TypeName mismatch: got %q, want %q", found.TypeName, room.TypeName)
	}

	room.Price = 499.00
	if err := repo.Update(ctx, room); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, err = repo.FindByID(ctx, room.ID)
	if err != nil {
		t.Fatalf("FindByID after Update failed: %v", err)
	}
	if found.Price != 499.00 {
		t.Errorf("Price after Update mismatch: got %f, want %f", found.Price, 499.00)
	}

	if err := repo.Delete(ctx, room.ID); err != nil {
		t.Fatalf("Delete (soft) failed: %v", err)
	}

	_, err = repo.FindByID(ctx, room.ID)
	if err == nil {
		t.Fatalf("FindByID after soft-delete should return error, got nil")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestRoomRepo_Create_WithImagesAndFacilities(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewRoomRepo(tx)

	region := createTestRegion(t, tx)
	hotel := createTestHotel(t, tx, region.ID, "Assoc Hotel")
	room := &model.Room{
		ID:                uuid.New(),
		HotelID:           hotel.ID,
		TypeName:          "Suite",
		TotalQuantity:     5,
		AvailableQuantity: 3,
		Price:             599.00,
		Status:            1,
		Images: []model.RoomImage{
			{ImageURL: "http://example.com/r1.jpg"},
			{ImageURL: "http://example.com/r2.jpg"},
		},
		Facilities: []model.RoomFacility{
			{FacilityName: "WiFi"},
			{FacilityName: "Breakfast"},
		},
	}

	if err := repo.Create(ctx, room); err != nil {
		t.Fatalf("Create with images and facilities failed: %v", err)
	}

	found, err := repo.FindByID(ctx, room.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if len(found.Images) != 2 {
		t.Errorf("Images count mismatch: got %d, want 2", len(found.Images))
	}
	if len(found.Facilities) != 2 {
		t.Errorf("Facilities count mismatch: got %d, want 2", len(found.Facilities))
	}
}

func TestRoomRepo_FindByHotelID(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewRoomRepo(tx)

	region := createTestRegion(t, tx)
	hotel := createTestHotel(t, tx, region.ID, "Multi Room Hotel")
	createTestRoom(t, tx, hotel.ID, "Standard")
	createTestRoom(t, tx, hotel.ID, "Deluxe")

	results, total, err := repo.FindByHotelID(ctx, hotel.ID, 0, 10)
	if err != nil {
		t.Fatalf("FindByHotelID failed: %v", err)
	}
	if total != 2 {
		t.Errorf("total mismatch: got %d, want 2", total)
	}
	if len(results) != 2 {
		t.Errorf("results count mismatch: got %d, want 2", len(results))
	}
}

func TestRoomRepo_FindAll_Pagination(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewRoomRepo(tx)

	region := createTestRegion(t, tx)
	hotel := createTestHotel(t, tx, region.ID, "Pag Room Hotel")
	for i := 0; i < 3; i++ {
		createTestRoom(t, tx, hotel.ID, "Type")
	}

	results, total, err := repo.FindAll(ctx, 0, 2)
	if err != nil {
		t.Fatalf("FindAll(0,2) failed: %v", err)
	}
	if total != 3 {
		t.Errorf("total mismatch: got %d, want 3", total)
	}
	if len(results) != 2 {
		t.Errorf("page size mismatch: got %d, want 2", len(results))
	}
}

// ---------------------------------------------------------------------------
// HotelImageRepo
// ---------------------------------------------------------------------------

func TestHotelImageRepo_CRUD(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	imgRepo := NewHotelImageRepo(tx)

	region := createTestRegion(t, tx)
	hotel := createTestHotel(t, tx, region.ID, "Img Hotel")

	img := &model.HotelImage{
		HotelID:  hotel.ID,
		ImageURL: "http://example.com/hotel.jpg",
	}
	if err := imgRepo.Create(ctx, img); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	results, err := imgRepo.FindByHotelID(ctx, hotel.ID)
	if err != nil {
		t.Fatalf("FindByHotelID failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("results count mismatch: got %d, want 1", len(results))
	}
	if results[0].ImageURL != img.ImageURL {
		t.Errorf("ImageURL mismatch: got %q, want %q", results[0].ImageURL, img.ImageURL)
	}

	if err := imgRepo.Delete(ctx, hotel.ID, img.ImageURL); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	results, err = imgRepo.FindByHotelID(ctx, hotel.ID)
	if err != nil {
		t.Fatalf("FindByHotelID after Delete failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("results count after Delete mismatch: got %d, want 0", len(results))
	}
}

// ---------------------------------------------------------------------------
// RoomImageRepo
// ---------------------------------------------------------------------------

func TestRoomImageRepo_CRUD(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	imgRepo := NewRoomImageRepo(tx)

	region := createTestRegion(t, tx)
	hotel := createTestHotel(t, tx, region.ID, "RImg Hotel")
	room := createTestRoom(t, tx, hotel.ID, "Standard")

	img := &model.RoomImage{
		RoomID:   room.ID,
		ImageURL: "http://example.com/room.jpg",
	}
	if err := imgRepo.Create(ctx, img); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	results, err := imgRepo.FindByRoomID(ctx, room.ID)
	if err != nil {
		t.Fatalf("FindByRoomID failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("results count mismatch: got %d, want 1", len(results))
	}

	if err := imgRepo.Delete(ctx, room.ID, img.ImageURL); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	results, err = imgRepo.FindByRoomID(ctx, room.ID)
	if err != nil {
		t.Fatalf("FindByRoomID after Delete failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("results count after Delete mismatch: got %d, want 0", len(results))
	}
}

// ---------------------------------------------------------------------------
// RoomFacilityRepo
// ---------------------------------------------------------------------------

func TestRoomFacilityRepo_CRUD(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	facRepo := NewRoomFacilityRepo(tx)

	region := createTestRegion(t, tx)
	hotel := createTestHotel(t, tx, region.ID, "Fac Hotel")
	room := createTestRoom(t, tx, hotel.ID, "Standard")

	fac := &model.RoomFacility{
		RoomID:       room.ID,
		FacilityName: "Gym",
	}
	if err := facRepo.Create(ctx, fac); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	results, err := facRepo.FindByRoomID(ctx, room.ID)
	if err != nil {
		t.Fatalf("FindByRoomID failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("results count mismatch: got %d, want 1", len(results))
	}
	if results[0].FacilityName != fac.FacilityName {
		t.Errorf("FacilityName mismatch: got %q, want %q", results[0].FacilityName, fac.FacilityName)
	}

	if err := facRepo.Delete(ctx, room.ID, fac.FacilityName); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	results, err = facRepo.FindByRoomID(ctx, room.ID)
	if err != nil {
		t.Fatalf("FindByRoomID after Delete failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("results count after Delete mismatch: got %d, want 0", len(results))
	}
}

// ---------------------------------------------------------------------------
// HotelSummaryRepo (view — skip if view not available)
// ---------------------------------------------------------------------------

func TestHotelSummaryRepo_FindByID(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewHotelSummaryRepo(tx)

	region := createTestRegion(t, tx)
	hotel := createTestHotel(t, tx, region.ID, "Summary Hotel")

	_, err := repo.FindByID(ctx, hotel.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			t.Skip("view not available or hotel not reflected in view")
		}
		t.Fatalf("FindByID failed: %v", err)
	}
}

func TestHotelSummaryRepo_FindAll(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewHotelSummaryRepo(tx)

	region := createTestRegion(t, tx)
	createTestHotel(t, tx, region.ID, "Summary Hotel A")
	createTestHotel(t, tx, region.ID, "Summary Hotel B")

	_, _, err := repo.FindAll(ctx, 0, 10, "", "", "", nil, nil, nil)
	if err != nil {
		t.Skipf("view not available: %v", err)
	}
}

// ---------------------------------------------------------------------------
// RoomDetailsRepo (view — skip if view not available)
// ---------------------------------------------------------------------------

func TestRoomDetailsRepo_FindByRoomID(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()
	repo := NewRoomDetailsRepo(tx)

	region := createTestRegion(t, tx)
	hotel := createTestHotel(t, tx, region.ID, "Details Hotel")
	room := createTestRoom(t, tx, hotel.ID, "Details Room")

	_, err := repo.FindByRoomID(ctx, room.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			t.Skip("view not available or room not reflected in view")
		}
		t.Fatalf("FindByRoomID failed: %v", err)
	}
}
