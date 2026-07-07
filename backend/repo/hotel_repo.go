package repo

import (
	"context"
	"time"

	model "backend/model/schema"
	"backend/model/view"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HotelRepo 酒店表（hotel_1718）的 GORM 仓库，支持软删除与硬删除。
type HotelRepo struct {
	db *gorm.DB
}

// NewHotelRepo 创建 HotelRepo 实例。
func NewHotelRepo(db *gorm.DB) *HotelRepo {
	return &HotelRepo{db: db}
}

// Create 插入酒店记录及其关联图片，使用事务保证原子性。
func (r *HotelRepo) Create(ctx context.Context, hotel *model.Hotel) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(hotel.Images) > 0 {
			for i := range hotel.Images {
				hotel.Images[i].HotelID = hotel.ID
			}
		}
		return tx.Create(hotel).Error
	})
}

// FindByID 根据酒店 ID 查询，预加载图片列表；若记录不存在返回 gorm.ErrRecordNotFound。
func (r *HotelRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Hotel, error) {
	var hotel model.Hotel
	err := r.db.WithContext(ctx).Preload("Images").Where("status != 0").First(&hotel, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &hotel, nil
}

// FindAll 查询酒店列表，支持按区域、星级、名称关键字筛选，返回分页结果与总记录数。
func (r *HotelRepo) FindAll(ctx context.Context, offset, limit int, regionID *int, starLevel *int16, keyword string) ([]model.Hotel, int64, error) {
	var results []model.Hotel
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Hotel{}).Where("status != 0")

	if regionID != nil && *regionID > 0 {
		query = query.Where("region_id = ?", *regionID)
	}
	if starLevel != nil && *starLevel > 0 {
		query = query.Where("star_level = ?", *starLevel)
	}
	if keyword != "" {
		query = query.Where("hotel_name ILIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Find(&results).Error
	return results, total, err
}

// Update 更新酒店记录（仅更新非零字段，不处理关联）。
func (r *HotelRepo) Update(ctx context.Context, hotel *model.Hotel) error {
	return r.db.WithContext(ctx).Save(hotel).Error
}

// Delete 软删除酒店记录，将 status 置为 0。
func (r *HotelRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.Hotel{}).Where("id = ?", id).Update("status", 0).Error
}

// HardDelete 硬删除酒店记录，并级联删除其关联图片。
func (r *HotelRepo) HardDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("hotel_id = ?", id).Delete(&model.HotelImage{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Hotel{}, "id = ?", id).Error
	})
}

// HotelImageRepo 酒店图片表（hotel_image_1718）的 GORM 仓库，使用硬删除。
type HotelImageRepo struct {
	db *gorm.DB
}

// NewHotelImageRepo 创建 HotelImageRepo 实例。
func NewHotelImageRepo(db *gorm.DB) *HotelImageRepo {
	return &HotelImageRepo{db: db}
}

// Create 插入一条酒店图片记录。
func (r *HotelImageRepo) Create(ctx context.Context, image *model.HotelImage) error {
	return r.db.WithContext(ctx).Create(image).Error
}

// FindByHotelID 根据酒店 ID 查询其全部图片记录。
func (r *HotelImageRepo) FindByHotelID(ctx context.Context, hotelID uuid.UUID) ([]model.HotelImage, error) {
	var results []model.HotelImage
	err := r.db.WithContext(ctx).Where("hotel_id = ?", hotelID).Find(&results).Error
	return results, err
}

// Delete 根据酒店 ID 与图片 URL 硬删除对应的酒店图片记录。
func (r *HotelImageRepo) Delete(ctx context.Context, hotelID uuid.UUID, imageURL string) error {
	return r.db.WithContext(ctx).Where("hotel_id = ? AND image_url = ?", hotelID, imageURL).Delete(&model.HotelImage{}).Error
}

// RoomRepo 客房表（room_1718）的 GORM 仓库，支持软删除与硬删除。
type RoomRepo struct {
	db *gorm.DB
}

// NewRoomRepo 创建 RoomRepo 实例。
func NewRoomRepo(db *gorm.DB) *RoomRepo {
	return &RoomRepo{db: db}
}

// Create 插入客房记录及其关联图片与设施，使用事务保证原子性。
func (r *RoomRepo) Create(ctx context.Context, room *model.Room) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(room.Images) > 0 {
			for i := range room.Images {
				room.Images[i].RoomID = room.ID
			}
		}
		if len(room.Facilities) > 0 {
			for i := range room.Facilities {
				room.Facilities[i].RoomID = room.ID
			}
		}
		return tx.Create(room).Error
	})
}

// FindByID 根据客房 ID 查询，预加载图片与设施列表；若记录不存在返回 gorm.ErrRecordNotFound。
func (r *RoomRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Room, error) {
	var room model.Room
	err := r.db.WithContext(ctx).Preload("Images").Preload("Facilities").Where("status != 0").First(&room, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// FindByHotelID 根据酒店 ID 查询其全部客房列表，返回分页结果与总记录数。
func (r *RoomRepo) FindByHotelID(ctx context.Context, hotelID uuid.UUID, offset, limit int) ([]model.Room, int64, error) {
	var results []model.Room
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Room{}).Where("hotel_id = ?", hotelID).Where("status != 0")
	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Preload("Images").Preload("Facilities").Find(&results).Error
	return results, total, err
}

// FindAvailableRooms 查询指定酒店在指定日期范围内仍有空余的客房列表。
// 筛选条件：available_quantity > 0 且不存在与该日期范围重叠的已确认订单将其完全订满。
func (r *RoomRepo) FindAvailableRooms(ctx context.Context, hotelID uuid.UUID, checkIn, checkOut time.Time, offset, limit int) ([]model.Room, int64, error) {
	var results []model.Room
	var total int64

	subQuery := r.db.WithContext(ctx).Model(&model.Order{}).
		Select("room_id").
		Where("status != ?", model.OrderCancelled).
		Where("check_in_date < ? AND check_out_date > ?", checkOut, checkIn)

	query := r.db.WithContext(ctx).Model(&model.Room{}).
		Where("hotel_id = ?", hotelID).
		Where("status != 0").
		Where("available_quantity > 0").
		Not("id IN (?)", subQuery)

	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Preload("Images").Preload("Facilities").Find(&results).Error
	return results, total, err
}

// FindAll 查询全部客房列表，返回分页结果与总记录数。
func (r *RoomRepo) FindAll(ctx context.Context, offset, limit int) ([]model.Room, int64, error) {
	var results []model.Room
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Room{}).Where("status != 0")
	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Preload("Images").Preload("Facilities").Find(&results).Error
	return results, total, err
}

// Update 更新客房记录（仅更新非零字段，不处理关联）。
func (r *RoomRepo) Update(ctx context.Context, room *model.Room) error {
	return r.db.WithContext(ctx).Save(room).Error
}

// Delete 软删除客房记录，将 status 置为 0。
func (r *RoomRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.Room{}).Where("id = ?", id).Update("status", 0).Error
}

// HardDelete 硬删除客房记录，并级联删除其关联图片与设施。
func (r *RoomRepo) HardDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("room_id = ?", id).Delete(&model.RoomImage{}).Error; err != nil {
			return err
		}
		if err := tx.Where("room_id = ?", id).Delete(&model.RoomFacility{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Room{}, "id = ?", id).Error
	})
}

// RoomImageRepo 客房图片表（room_image_1718）的 GORM 仓库，使用硬删除。
type RoomImageRepo struct {
	db *gorm.DB
}

// NewRoomImageRepo 创建 RoomImageRepo 实例。
func NewRoomImageRepo(db *gorm.DB) *RoomImageRepo {
	return &RoomImageRepo{db: db}
}

// Create 插入一条客房图片记录。
func (r *RoomImageRepo) Create(ctx context.Context, image *model.RoomImage) error {
	return r.db.WithContext(ctx).Create(image).Error
}

// FindByRoomID 根据客房 ID 查询其全部图片记录。
func (r *RoomImageRepo) FindByRoomID(ctx context.Context, roomID uuid.UUID) ([]model.RoomImage, error) {
	var results []model.RoomImage
	err := r.db.WithContext(ctx).Where("room_id = ?", roomID).Find(&results).Error
	return results, err
}

// Delete 根据客房 ID 与图片 URL 硬删除对应的客房图片记录。
func (r *RoomImageRepo) Delete(ctx context.Context, roomID uuid.UUID, imageURL string) error {
	return r.db.WithContext(ctx).Where("room_id = ? AND image_url = ?", roomID, imageURL).Delete(&model.RoomImage{}).Error
}

// RoomFacilityRepo 客房设施表（room_facility_1718）的 GORM 仓库，使用硬删除。
type RoomFacilityRepo struct {
	db *gorm.DB
}

// NewRoomFacilityRepo 创建 RoomFacilityRepo 实例。
func NewRoomFacilityRepo(db *gorm.DB) *RoomFacilityRepo {
	return &RoomFacilityRepo{db: db}
}

// Create 插入一条客房设施记录。
func (r *RoomFacilityRepo) Create(ctx context.Context, facility *model.RoomFacility) error {
	return r.db.WithContext(ctx).Create(facility).Error
}

// FindByRoomID 根据客房 ID 查询其全部设施记录。
func (r *RoomFacilityRepo) FindByRoomID(ctx context.Context, roomID uuid.UUID) ([]model.RoomFacility, error) {
	var results []model.RoomFacility
	err := r.db.WithContext(ctx).Where("room_id = ?", roomID).Find(&results).Error
	return results, err
}

// Delete 根据客房 ID 与设施名称硬删除对应的客房设施记录。
func (r *RoomFacilityRepo) Delete(ctx context.Context, roomID uuid.UUID, facilityName string) error {
	return r.db.WithContext(ctx).Where("room_id = ? AND facility_name = ?", roomID, facilityName).Delete(&model.RoomFacility{}).Error
}

// HotelSummaryRepo 酒店摘要视图（view_hotel_summary_1718）的只读仓库。
type HotelSummaryRepo struct {
	db *gorm.DB
}

// NewHotelSummaryRepo 创建 HotelSummaryRepo 实例。
func NewHotelSummaryRepo(db *gorm.DB) *HotelSummaryRepo {
	return &HotelSummaryRepo{db: db}
}

// FindByID 根据酒店 ID 查询视图记录。
func (r *HotelSummaryRepo) FindByID(ctx context.Context, hotelID uuid.UUID) (*view.HotelSummary, error) {
	var result view.HotelSummary
	err := r.db.WithContext(ctx).Where("hotel_id = ?", hotelID).Where("status = 1").First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FindAll 查询酒店摘要视图列表，支持按省、市、区、星级、价格范围筛选，返回分页结果与总记录数。
// 仅当参数非零时应用对应过滤条件；结果按 avg_rating 降序排列。
func (r *HotelSummaryRepo) FindAll(ctx context.Context, offset, limit int, province, city, district string, starLevel *int16, minPrice, maxPrice *float64) ([]view.HotelSummary, int64, error) {
	var results []view.HotelSummary
	var total int64
	query := r.db.WithContext(ctx).Model(&view.HotelSummary{}).Where("status = 1")

	if province != "" {
		query = query.Where("province = ?", province)
	}
	if city != "" {
		query = query.Where("city = ?", city)
	}
	if district != "" {
		query = query.Where("district = ?", district)
	}
	if starLevel != nil && *starLevel > 0 {
		query = query.Where("star_level = ?", *starLevel)
	}
	if minPrice != nil && *minPrice > 0 {
		query = query.Where("min_price >= ?", *minPrice)
	}
	if maxPrice != nil && *maxPrice > 0 {
		query = query.Where("min_price <= ?", *maxPrice)
	}

	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Order("avg_rating DESC").Find(&results).Error
	return results, total, err
}

// FindByRegionID 根据区域 ID 查询酒店摘要视图列表，返回分页结果与总记录数。
func (r *HotelSummaryRepo) FindByRegionID(ctx context.Context, regionID int, offset, limit int) ([]view.HotelSummary, int64, error) {
	var results []view.HotelSummary
	var total int64
	query := r.db.WithContext(ctx).Model(&view.HotelSummary{}).
		Where("region_id = ?", regionID).
		Where("status = 1")
	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Order("avg_rating DESC").Find(&results).Error
	return results, total, err
}

// RoomDetailsRepo 客房详情视图（view_room_details_1718）的只读仓库。
type RoomDetailsRepo struct {
	db *gorm.DB
}

// NewRoomDetailsRepo 创建 RoomDetailsRepo 实例。
func NewRoomDetailsRepo(db *gorm.DB) *RoomDetailsRepo {
	return &RoomDetailsRepo{db: db}
}

// FindByRoomID 根据客房 ID 查询视图记录。
func (r *RoomDetailsRepo) FindByRoomID(ctx context.Context, roomID uuid.UUID) (*view.RoomDetails, error) {
	var result view.RoomDetails
	err := r.db.WithContext(ctx).Where("room_id = ?", roomID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FindByHotelID 根据酒店 ID 查询客房详情视图列表，返回分页结果与总记录数。
func (r *RoomDetailsRepo) FindByHotelID(ctx context.Context, hotelID uuid.UUID, offset, limit int) ([]view.RoomDetails, int64, error) {
	var results []view.RoomDetails
	var total int64
	query := r.db.WithContext(ctx).Model(&view.RoomDetails{}).Where("hotel_id = ?", hotelID)
	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Order("hotel_name, price").Find(&results).Error
	return results, total, err
}

// FindAll 查询客房详情视图列表，支持按省、市、区、星级、价格范围筛选，返回分页结果与总记录数。
// 仅当参数非零时应用对应过滤条件；结果按 hotel_name、price 排序。
func (r *RoomDetailsRepo) FindAll(ctx context.Context, offset, limit int, province, city, district string, starLevel *int16, minPrice, maxPrice *float64) ([]view.RoomDetails, int64, error) {
	var results []view.RoomDetails
	var total int64
	query := r.db.WithContext(ctx).Model(&view.RoomDetails{})

	if province != "" {
		query = query.Where("province = ?", province)
	}
	if city != "" {
		query = query.Where("city = ?", city)
	}
	if district != "" {
		query = query.Where("district = ?", district)
	}
	if starLevel != nil && *starLevel > 0 {
		query = query.Where("star_level = ?", *starLevel)
	}
	if minPrice != nil && *minPrice > 0 {
		query = query.Where("price >= ?", *minPrice)
	}
	if maxPrice != nil && *maxPrice > 0 {
		query = query.Where("price <= ?", *maxPrice)
	}

	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Order("hotel_name, price").Find(&results).Error
	return results, total, err
}
