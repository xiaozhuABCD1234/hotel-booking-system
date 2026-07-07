package repo

import (
	"context"
	"errors"

	model "backend/model/schema"
	"backend/model/view"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepo 用户表（user_1718）的 GORM 仓库，支持软删除。
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo 创建 UserRepo 实例。
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create 插入一条用户记录。
func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID 根据用户 ID 查询，预加载 VipLevel，排除已软删除用户。
func (r *UserRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("VipLevel").Where("status != ?", 0).First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名精确查询，排除已软删除用户。
func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ? AND status != ?", username, 0).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByPhone 根据手机号精确查询，排除已软删除用户。
func (r *UserRepo) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("phone = ? AND status != ?", phone, 0).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindAll 查询用户列表，支持按角色筛选、分页，预加载 VipLevel，按 create_at 降序，排除已软删除用户。
func (r *UserRepo) FindAll(ctx context.Context, offset, limit int, role *model.UserRole) ([]model.User, int64, error) {
	var results []model.User
	var total int64
	query := r.db.WithContext(ctx).Model(&model.User{}).Where("status != ?", 0)
	if role != nil {
		query = query.Where("role = ?", *role)
	}
	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Preload("VipLevel").Order("create_at DESC").Find(&results).Error
	return results, total, err
}

// Update 更新用户非零字段，不修改 CreateAt。
func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(user).Error
}

// UpdatePassword 更新用户密码。
func (r *UserRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}

// UpdatePoints 更新用户积分。
func (r *UserRepo) UpdatePoints(ctx context.Context, userID uuid.UUID, points int32) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("points", points).Error
}

// UpdateVipLevel 更新用户 VIP 等级。
func (r *UserRepo) UpdateVipLevel(ctx context.Context, userID uuid.UUID, vipLevel int16) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("vip_level", vipLevel).Error
}

// Delete 软删除用户，将 status 置为 0。
func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("status", 0).Error
}

// HardDelete 根据 ID 硬删除用户记录。
func (r *UserRepo) HardDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}

// VipLevelRepo VIP 等级表（vip_level_1718）的 GORM 仓库，支持硬删除。
type VipLevelRepo struct {
	db *gorm.DB
}

// NewVipLevelRepo 创建 VipLevelRepo 实例。
func NewVipLevelRepo(db *gorm.DB) *VipLevelRepo {
	return &VipLevelRepo{db: db}
}

// FindByLevel 根据等级查询 VIP 定义。
func (r *VipLevelRepo) FindByLevel(ctx context.Context, level int16) (*model.VipLevel, error) {
	var vip model.VipLevel
	err := r.db.WithContext(ctx).First(&vip, "level = ?", level).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &vip, nil
}

// FindAll 查询全部 VIP 等级定义，按 level 升序。
func (r *VipLevelRepo) FindAll(ctx context.Context) ([]model.VipLevel, error) {
	var results []model.VipLevel
	err := r.db.WithContext(ctx).Order("level ASC").Find(&results).Error
	return results, err
}

// FindByMinPoints 根据积分查询适用的最高 VIP 等级（min_points <= points，按 level 降序取第一条）。
func (r *VipLevelRepo) FindByMinPoints(ctx context.Context, points int32) (*model.VipLevel, error) {
	var vip model.VipLevel
	err := r.db.WithContext(ctx).Where("min_points <= ?", points).Order("level DESC").First(&vip).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &vip, nil
}

// Create 插入一条 VIP 等级定义。
func (r *VipLevelRepo) Create(ctx context.Context, vip *model.VipLevel) error {
	return r.db.WithContext(ctx).Create(vip).Error
}

// Update 更新 VIP 等级定义。
func (r *VipLevelRepo) Update(ctx context.Context, vip *model.VipLevel) error {
	return r.db.WithContext(ctx).Save(vip).Error
}

// Delete 根据等级硬删除 VIP 定义。
func (r *VipLevelRepo) Delete(ctx context.Context, level int16) error {
	return r.db.WithContext(ctx).Delete(&model.VipLevel{}, "level = ?", level).Error
}

// UserVipRepo 用户 VIP 视图（view_user_vip_1718）的只读仓库。
type UserVipRepo struct {
	db *gorm.DB
}

// NewUserVipRepo 创建 UserVipRepo 实例。
func NewUserVipRepo(db *gorm.DB) *UserVipRepo {
	return &UserVipRepo{db: db}
}

// FindByUserID 根据用户 ID 查询视图记录。
func (r *UserVipRepo) FindByUserID(ctx context.Context, userID uuid.UUID) (*view.UserVip, error) {
	var result view.UserVip
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &result, nil
}

// FindAll 查询全部用户 VIP 视图，支持按角色筛选、分页，按 points 降序。
func (r *UserVipRepo) FindAll(ctx context.Context, offset, limit int, role string) ([]view.UserVip, int64, error) {
	var results []view.UserVip
	var total int64
	query := r.db.WithContext(ctx).Model(&view.UserVip{})
	if role != "" {
		query = query.Where("role = ?", role)
	}
	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Order("points DESC").Find(&results).Error
	return results, total, err
}

// FindByVipLevel 根据 VIP 等级查询视图列表，支持分页，按 points 降序。
func (r *UserVipRepo) FindByVipLevel(ctx context.Context, level int16, offset, limit int) ([]view.UserVip, int64, error) {
	var results []view.UserVip
	var total int64
	query := r.db.WithContext(ctx).Model(&view.UserVip{}).Where("vip_level = ?", level)
	query.Count(&total)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Order("points DESC").Find(&results).Error
	return results, total, err
}
