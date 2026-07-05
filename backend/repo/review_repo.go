package repo

import (
	"context"

	model "backend/model/schema"
	"backend/model/view"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReviewRepo 评价表（review_1718）的 GORM 仓库，支持硬删除。
type ReviewRepo struct {
	db *gorm.DB
}

// NewReviewRepo 创建 ReviewRepo 实例。
func NewReviewRepo(db *gorm.DB) *ReviewRepo {
	return &ReviewRepo{db: db}
}

// Create 插入一条评价记录。
func (r *ReviewRepo) Create(ctx context.Context, review *model.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

// FindByID 根据评价 ID 查询，预加载 User、Hotel、Order。
func (r *ReviewRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Review, error) {
	var review model.Review
	err := r.db.WithContext(ctx).Preload("User").Preload("Hotel").Preload("Order").First(&review, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// FindByHotelID 根据酒店 ID 查询评价列表，预加载 User，按 create_at 降序。
func (r *ReviewRepo) FindByHotelID(ctx context.Context, hotelID uuid.UUID, offset, limit int) ([]model.Review, int64, error) {
	var results []model.Review
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Review{}).Where("hotel_id = ?", hotelID)
	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Preload("User").Order("create_at DESC").Find(&results).Error
	return results, total, err
}

// FindByUserID 根据用户 ID 查询评价列表，预加载 Hotel，按 create_at 降序。
func (r *ReviewRepo) FindByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]model.Review, int64, error) {
	var results []model.Review
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Review{}).Where("user_id = ?", userID)
	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Preload("Hotel").Order("create_at DESC").Find(&results).Error
	return results, total, err
}

// FindByOrderID 根据订单 ID 查询唯一评价（user_id, order_id 唯一约束）。
func (r *ReviewRepo) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*model.Review, error) {
	var review model.Review
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// FindAll 查询全部评价，预加载 User、Hotel，按 create_at 降序。
func (r *ReviewRepo) FindAll(ctx context.Context, offset, limit int) ([]model.Review, int64, error) {
	var results []model.Review
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Review{})
	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Preload("User").Preload("Hotel").Order("create_at DESC").Find(&results).Error
	return results, total, err
}

// FindByRating 根据评分查询评价列表，按 create_at 降序。
func (r *ReviewRepo) FindByRating(ctx context.Context, rating int16, offset, limit int) ([]model.Review, int64, error) {
	var results []model.Review
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Review{}).Where("rating = ?", rating)
	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Order("create_at DESC").Find(&results).Error
	return results, total, err
}

// Update 更新评价的 Rating 和 Content，不允许修改 UserID/OrderID/HotelID。
func (r *ReviewRepo) Update(ctx context.Context, review *model.Review) error {
	return r.db.WithContext(ctx).Model(review).Select("Rating", "Content").Updates(review).Error
}

// Delete 根据 ID 硬删除评价记录。
func (r *ReviewRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Review{}, "id = ?", id).Error
}

// ReviewFullRepo 评价视图（view_review_full_1718）的只读仓库。
type ReviewFullRepo struct {
	db *gorm.DB
}

// NewReviewFullRepo 创建 ReviewFullRepo 实例。
func NewReviewFullRepo(db *gorm.DB) *ReviewFullRepo {
	return &ReviewFullRepo{db: db}
}

// FindByReviewID 根据评价 ID 查询视图记录。
func (r *ReviewFullRepo) FindByReviewID(ctx context.Context, reviewID uuid.UUID) (*view.ReviewFull, error) {
	var result view.ReviewFull
	err := r.db.WithContext(ctx).Where("review_id = ?", reviewID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FindByHotelID 根据酒店 ID 查询视图列表，按 create_at 降序。
func (r *ReviewFullRepo) FindByHotelID(ctx context.Context, hotelID uuid.UUID, offset, limit int) ([]view.ReviewFull, int64, error) {
	var results []view.ReviewFull
	var total int64
	query := r.db.WithContext(ctx).Model(&view.ReviewFull{}).Where("hotel_id = ?", hotelID)
	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Order("create_at DESC").Find(&results).Error
	return results, total, err
}

// FindByUserID 根据用户 ID 查询视图列表，按 create_at 降序。
func (r *ReviewFullRepo) FindByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]view.ReviewFull, int64, error) {
	var results []view.ReviewFull
	var total int64
	query := r.db.WithContext(ctx).Model(&view.ReviewFull{}).Where("user_id = ?", userID)
	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Order("create_at DESC").Find(&results).Error
	return results, total, err
}

// FindAll 查询全部视图记录，按 create_at 降序。
func (r *ReviewFullRepo) FindAll(ctx context.Context, offset, limit int) ([]view.ReviewFull, int64, error) {
	var results []view.ReviewFull
	var total int64
	query := r.db.WithContext(ctx).Model(&view.ReviewFull{})
	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Order("create_at DESC").Find(&results).Error
	return results, total, err
}
