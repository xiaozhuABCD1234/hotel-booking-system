package repo

import (
	"context"

	model "backend/model/schema"
	"backend/model/view"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ===================== PersonRepo =====================

type PersonRepo struct {
	db *gorm.DB
}

func NewPersonRepo(db *gorm.DB) *PersonRepo {
	return &PersonRepo{db: db}
}

func (r *PersonRepo) Create(ctx context.Context, person *model.Person) error {
	return r.db.WithContext(ctx).Create(person).Error
}

func (r *PersonRepo) FindByIDCard(ctx context.Context, idCard string) (*model.Person, error) {
	var person model.Person
	if err := r.db.WithContext(ctx).Where("id_card = ?", idCard).First(&person).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *PersonRepo) FindAll(ctx context.Context, offset, limit int, keyword string) ([]model.Person, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Person{})
	if keyword != "" {
		query = query.Where("name ILIKE ?", "%"+keyword+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("name")
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	var persons []model.Person
	if err := query.Find(&persons).Error; err != nil {
		return nil, 0, err
	}
	return persons, total, nil
}

func (r *PersonRepo) Update(ctx context.Context, person *model.Person) error {
	return r.db.WithContext(ctx).Select("name", "phone").Where("id_card = ?", person.IDCard).Updates(person).Error
}

func (r *PersonRepo) Delete(ctx context.Context, idCard string) error {
	return r.db.WithContext(ctx).Where("id_card = ?", idCard).Delete(&model.Person{}).Error
}

func (r *PersonRepo) Upsert(ctx context.Context, person *model.Person) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id_card"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "phone"}),
	}).Create(person).Error
}

// ===================== RegionRepo =====================

type RegionRepo struct {
	db *gorm.DB
}

func NewRegionRepo(db *gorm.DB) *RegionRepo {
	return &RegionRepo{db: db}
}

func (r *RegionRepo) FindByID(ctx context.Context, id int) (*model.Region, error) {
	var region model.Region
	if err := r.db.WithContext(ctx).Preload("Parent").Where("id = ?", id).First(&region).Error; err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *RegionRepo) FindByParentID(ctx context.Context, parentID int) ([]model.Region, error) {
	var regions []model.Region
	if err := r.db.WithContext(ctx).Where("parents_id = ?", parentID).Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

func (r *RegionRepo) FindAllProvinces(ctx context.Context) ([]model.Region, error) {
	var regions []model.Region
	if err := r.db.WithContext(ctx).Where("parents_id IS NULL").Order("id").Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

func (r *RegionRepo) FindAll(ctx context.Context) ([]model.Region, error) {
	var regions []model.Region
	if err := r.db.WithContext(ctx).Preload("Parent").Order("id").Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

func (r *RegionRepo) Create(ctx context.Context, region *model.Region) error {
	return r.db.WithContext(ctx).Create(region).Error
}

func (r *RegionRepo) Update(ctx context.Context, region *model.Region) error {
	return r.db.WithContext(ctx).Save(region).Error
}

func (r *RegionRepo) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Region{}).Error
}

// ===================== PersonInfoRepo =====================

type PersonInfoRepo struct {
	db *gorm.DB
}

func NewPersonInfoRepo(db *gorm.DB) *PersonInfoRepo {
	return &PersonInfoRepo{db: db}
}

func (r *PersonInfoRepo) FindByIDCard(ctx context.Context, idCard string) (*view.PersonInfo, error) {
	var info view.PersonInfo
	if err := r.db.WithContext(ctx).Where("id_card = ?", idCard).First(&info).Error; err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *PersonInfoRepo) FindAll(ctx context.Context, offset, limit int, gender string, minAge, maxAge *int) ([]view.PersonInfo, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&view.PersonInfo{})
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}
	if minAge != nil {
		query = query.Where("age >= ?", *minAge)
	}
	if maxAge != nil {
		query = query.Where("age <= ?", *maxAge)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("name")
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	var infos []view.PersonInfo
	if err := query.Find(&infos).Error; err != nil {
		return nil, 0, err
	}
	return infos, total, nil
}

// ===================== GuestBookingStatsRepo =====================

type GuestBookingStatsRepo struct {
	db *gorm.DB
}

func NewGuestBookingStatsRepo(db *gorm.DB) *GuestBookingStatsRepo {
	return &GuestBookingStatsRepo{db: db}
}

func (r *GuestBookingStatsRepo) FindByIDCard(ctx context.Context, idCard string) (*view.GuestBookingStats, error) {
	var stats view.GuestBookingStats
	if err := r.db.WithContext(ctx).Where("person_id_card = ?", idCard).First(&stats).Error; err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *GuestBookingStatsRepo) FindAll(ctx context.Context, offset, limit int, ageGroup string, gender string, favCity string) ([]view.GuestBookingStats, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&view.GuestBookingStats{})
	if ageGroup != "" {
		query = query.Where("age_group ILIKE ?", "%"+ageGroup+"%")
	}
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}
	if favCity != "" {
		query = query.Where("fav_city ILIKE ?", "%"+favCity+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("total_amount DESC")
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	var stats []view.GuestBookingStats
	if err := query.Find(&stats).Error; err != nil {
		return nil, 0, err
	}
	return stats, total, nil
}

func (r *GuestBookingStatsRepo) FindTopGuests(ctx context.Context, limit int) ([]view.GuestBookingStats, error) {
	var stats []view.GuestBookingStats
	if err := r.db.WithContext(ctx).Order("total_amount DESC").Limit(limit).Find(&stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}
