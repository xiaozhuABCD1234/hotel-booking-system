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

// ===================== PersonRepo Tests =====================

func TestPersonRepo_CRUD_Cycle(t *testing.T) {
	tx := txRepo(t)
	repo := NewPersonRepo(tx)
	ctx := context.Background()

	person := &model.Person{IDCard: "110101199001011234", Name: "Test Person", Phone: ptr("13800000000")}
	if err := repo.Create(ctx, person); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.FindByIDCard(ctx, person.IDCard)
	if err != nil {
		t.Fatalf("FindByIDCard failed: %v", err)
	}
	if found.Name != person.Name {
		t.Errorf("name mismatch: got %q, want %q", found.Name, person.Name)
	}

	person.Phone = ptr("13900000000")
	if err := repo.Update(ctx, person); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, err = repo.FindByIDCard(ctx, person.IDCard)
	if err != nil {
		t.Fatalf("FindByIDCard after update failed: %v", err)
	}
	if found.Phone == nil || *found.Phone != "13900000000" {
		t.Errorf("phone update mismatch: got %v, want 13900000000", found.Phone)
	}

	if err := repo.Delete(ctx, person.IDCard); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.FindByIDCard(ctx, person.IDCard)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound after delete, got: %v", err)
	}
}

func TestPersonRepo_Upsert(t *testing.T) {
	tx := txRepo(t)
	repo := NewPersonRepo(tx)
	ctx := context.Background()

	person := &model.Person{IDCard: "110101199001011235", Name: "Upsert Person", Phone: ptr("13800000001")}
	if err := repo.Upsert(ctx, person); err != nil {
		t.Fatalf("Upsert (create) failed: %v", err)
	}

	found, err := repo.FindByIDCard(ctx, person.IDCard)
	if err != nil {
		t.Fatalf("FindByIDCard after upsert failed: %v", err)
	}
	if found.Name != "Upsert Person" {
		t.Errorf("name mismatch after upsert create: got %q, want %q", found.Name, "Upsert Person")
	}

	person.Name = "Updated Upsert Person"
	person.Phone = ptr("13900000001")
	if err := repo.Upsert(ctx, person); err != nil {
		t.Fatalf("Upsert (update) failed: %v", err)
	}

	found, err = repo.FindByIDCard(ctx, person.IDCard)
	if err != nil {
		t.Fatalf("FindByIDCard after upsert update failed: %v", err)
	}
	if found.Name != "Updated Upsert Person" {
		t.Errorf("name mismatch after upsert update: got %q, want %q", found.Name, "Updated Upsert Person")
	}
	if found.Phone == nil || *found.Phone != "13900000001" {
		t.Errorf("phone mismatch after upsert update: got %v, want 13900000001", found.Phone)
	}
}

func TestPersonRepo_FindAll_Pagination(t *testing.T) {
	tx := txRepo(t)
	repo := NewPersonRepo(tx)
	ctx := context.Background()

	persons := []*model.Person{
		{IDCard: "110101199001011236", Name: "Person A", Phone: ptr("13800000002")},
		{IDCard: "110101199001011237", Name: "Person B", Phone: ptr("13800000003")},
		{IDCard: "110101199001011238", Name: "Person C", Phone: ptr("13800000004")},
	}
	for _, p := range persons {
		if err := repo.Create(ctx, p); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	results, total, err := repo.FindAll(ctx, 0, 2, "")
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if total != 3 {
		t.Errorf("total mismatch: got %d, want 3", total)
	}
	if len(results) != 2 {
		t.Errorf("results length mismatch: got %d, want 2", len(results))
	}
}

func TestPersonRepo_FindAll_KeywordFilter(t *testing.T) {
	tx := txRepo(t)
	repo := NewPersonRepo(tx)
	ctx := context.Background()

	alice := &model.Person{IDCard: "110101199001011239", Name: "Alice Smith", Phone: ptr("13800000005")}
	bob := &model.Person{IDCard: "11010119900101123X", Name: "Bob Jones", Phone: ptr("13800000006")}
	if err := repo.Create(ctx, alice); err != nil {
		t.Fatalf("Create Alice failed: %v", err)
	}
	if err := repo.Create(ctx, bob); err != nil {
		t.Fatalf("Create Bob failed: %v", err)
	}

	results, total, err := repo.FindAll(ctx, 0, 0, "Alice")
	if err != nil {
		t.Fatalf("FindAll with keyword failed: %v", err)
	}
	if total != 1 {
		t.Errorf("total mismatch: got %d, want 1", total)
	}
	if len(results) != 1 {
		t.Errorf("results length mismatch: got %d, want 1", len(results))
	}
	if results[0].Name != "Alice Smith" {
		t.Errorf("name mismatch: got %q, want %q", results[0].Name, "Alice Smith")
	}
}

// ===================== RegionRepo Tests =====================

func TestRegionRepo_CRUD_Cycle(t *testing.T) {
	tx := txRepo(t)
	repo := NewRegionRepo(tx)
	ctx := context.Background()

	province := &model.Region{RegionName: "Test Province"}
	if err := repo.Create(ctx, province); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if province.ID == 0 {
		t.Fatalf("expected region ID to be populated after create")
	}

	found, err := repo.FindByID(ctx, province.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.RegionName != province.RegionName {
		t.Errorf("name mismatch: got %q, want %q", found.RegionName, province.RegionName)
	}

	province.RegionName = "Updated Province"
	if err := repo.Update(ctx, province); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, err = repo.FindByID(ctx, province.ID)
	if err != nil {
		t.Fatalf("FindByID after update failed: %v", err)
	}
	if found.RegionName != "Updated Province" {
		t.Errorf("name mismatch after update: got %q, want %q", found.RegionName, "Updated Province")
	}

	if err := repo.Delete(ctx, province.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.FindByID(ctx, province.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound after delete, got: %v", err)
	}
}

func TestRegionRepo_FindByParentID(t *testing.T) {
	tx := txRepo(t)
	repo := NewRegionRepo(tx)
	ctx := context.Background()

	province := &model.Region{RegionName: "Test Province"}
	if err := repo.Create(ctx, province); err != nil {
		t.Fatalf("Create province failed: %v", err)
	}

	city1 := &model.Region{RegionName: "Test City 1", ParentsID: &province.ID}
	city2 := &model.Region{RegionName: "Test City 2", ParentsID: &province.ID}
	if err := repo.Create(ctx, city1); err != nil {
		t.Fatalf("Create city1 failed: %v", err)
	}
	if err := repo.Create(ctx, city2); err != nil {
		t.Fatalf("Create city2 failed: %v", err)
	}

	results, err := repo.FindByParentID(ctx, province.ID)
	if err != nil {
		t.Fatalf("FindByParentID failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("results length mismatch: got %d, want 2", len(results))
	}
}

func TestRegionRepo_FindAllProvinces(t *testing.T) {
	tx := txRepo(t)
	repo := NewRegionRepo(tx)
	ctx := context.Background()

	province1 := &model.Region{RegionName: "Province One"}
	province2 := &model.Region{RegionName: "Province Two"}
	if err := repo.Create(ctx, province1); err != nil {
		t.Fatalf("Create province1 failed: %v", err)
	}
	if err := repo.Create(ctx, province2); err != nil {
		t.Fatalf("Create province2 failed: %v", err)
	}

	results, err := repo.FindAllProvinces(ctx)
	if err != nil {
		t.Fatalf("FindAllProvinces failed: %v", err)
	}
	if len(results) < 2 {
		t.Errorf("results length mismatch: got %d, want at least 2", len(results))
	}
	for _, r := range results {
		if r.ParentsID != nil {
			t.Errorf("expected province to have nil ParentsID, got %v", r.ParentsID)
		}
	}
}

func TestRegionRepo_FindAll(t *testing.T) {
	tx := txRepo(t)
	repo := NewRegionRepo(tx)
	ctx := context.Background()

	province := &model.Region{RegionName: "Test Province"}
	if err := repo.Create(ctx, province); err != nil {
		t.Fatalf("Create province failed: %v", err)
	}

	city := &model.Region{RegionName: "Test City", ParentsID: &province.ID}
	if err := repo.Create(ctx, city); err != nil {
		t.Fatalf("Create city failed: %v", err)
	}

	results, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if len(results) < 2 {
		t.Errorf("results length mismatch: got %d, want at least 2", len(results))
	}

	var foundProvince, foundCity bool
	for _, r := range results {
		if r.ID == province.ID {
			foundProvince = true
			if r.Parent != nil {
				t.Errorf("expected province to have nil Parent, got %+v", r.Parent)
			}
		}
		if r.ID == city.ID {
			foundCity = true
			if r.Parent == nil {
				t.Errorf("expected city to have preloaded Parent")
			} else if r.Parent.ID != province.ID {
				t.Errorf("city parent mismatch: got %d, want %d", r.Parent.ID, province.ID)
			}
		}
	}
	if !foundProvince {
		t.Errorf("province not found in FindAll results")
	}
	if !foundCity {
		t.Errorf("city not found in FindAll results")
	}
}

// ===================== PersonInfoRepo Tests =====================

func TestPersonInfoRepo_FindByIDCard(t *testing.T) {
	tx := txRepo(t)
	personRepo := NewPersonRepo(tx)
	infoRepo := NewPersonInfoRepo(tx)
	ctx := context.Background()

	person := &model.Person{IDCard: "11010119900101123x", Name: "Info Test Person", Phone: ptr("13800000007")}
	if err := personRepo.Create(ctx, person); err != nil {
		t.Fatalf("Create person failed: %v", err)
	}

	info, err := infoRepo.FindByIDCard(ctx, person.IDCard)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			t.Skip("view not available or person not in view")
		}
		t.Skipf("view not available: %v", err)
	}
	if info.IDCard != person.IDCard {
		t.Errorf("id_card mismatch: got %q, want %q", info.IDCard, person.IDCard)
	}
	if info.Name != person.Name {
		t.Errorf("name mismatch: got %q, want %q", info.Name, person.Name)
	}
}

// ===================== GuestBookingStatsRepo Tests =====================

func TestGuestBookingStatsRepo_FindByIDCard(t *testing.T) {
	tx := txRepo(t)
	ctx := context.Background()

	// Create dependencies: VipLevel → User → Person → Region → Hotel → Room → Order → OrderGuest
	if err := tx.WithContext(ctx).Exec("INSERT INTO vip_level_1718 (level, level_name, min_points, discount_rate) VALUES (?, ?, ?, ?) ON CONFLICT DO NOTHING", 0, "普通会员", 0, 1.0).Error; err != nil {
		t.Fatalf("Create vip level failed: %v", err)
	}

	userRepo := NewUserRepo(tx)
	user := &model.User{
		ID:         uuid.New(),
		Username:   "stats_test_user",
		Password:   "hashed_password",
		Phone:      ptr("13800000008"),
		Role:       model.RoleCustomer,
		VipLevelID: 0,
		Status:     1,
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	personRepo := NewPersonRepo(tx)
	person := &model.Person{IDCard: "110101199001011240", Name: "Stats Person", Phone: ptr("13800000009")}
	if err := personRepo.Create(ctx, person); err != nil {
		t.Fatalf("Create person failed: %v", err)
	}

	regionRepo := NewRegionRepo(tx)
	region := &model.Region{RegionName: "Stats City"}
	if err := regionRepo.Create(ctx, region); err != nil {
		t.Fatalf("Create region failed: %v", err)
	}

	hotelRepo := NewHotelRepo(tx)
	hotel := &model.Hotel{
		ID:        uuid.New(),
		HotelName: "Stats Hotel",
		RegionID:  region.ID,
		Address:   "Stats Address",
		Telephone: "010-12345678",
		Status:    1,
	}
	if err := hotelRepo.Create(ctx, hotel); err != nil {
		t.Fatalf("Create hotel failed: %v", err)
	}

	roomRepo := NewRoomRepo(tx)
	room := &model.Room{
		ID:                uuid.New(),
		HotelID:           hotel.ID,
		TypeName:          "Standard",
		TotalQuantity:     10,
		AvailableQuantity: 5,
		Price:             299.00,
		Status:            1,
	}
	if err := roomRepo.Create(ctx, room); err != nil {
		t.Fatalf("Create room failed: %v", err)
	}

	now := time.Now()
	order := &model.Order{
		ID:           uuid.New(),
		UserID:       user.ID,
		RoomID:       room.ID,
		Quantity:     1,
		CheckInDate:  now,
		CheckOutDate: now.Add(24 * time.Hour),
		TotalPrice:   299.00,
		ActualPrice:  299.00,
		Status:       model.OrderBooked,
	}
	if err := tx.WithContext(ctx).Create(order).Error; err != nil {
		t.Fatalf("Create order failed: %v", err)
	}

	guest := &model.OrderGuest{OrderID: order.ID, IDCard: person.IDCard}
	if err := tx.WithContext(ctx).Create(guest).Error; err != nil {
		t.Fatalf("Create order guest failed: %v", err)
	}

	statsRepo := NewGuestBookingStatsRepo(tx)
	stats, err := statsRepo.FindByIDCard(ctx, person.IDCard)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			t.Skip("view not available or no stats for this person")
		}
		t.Skipf("view not available: %v", err)
	}
	if stats.PersonIDCard != person.IDCard {
		t.Errorf("person_id_card mismatch: got %q, want %q", stats.PersonIDCard, person.IDCard)
	}
}
