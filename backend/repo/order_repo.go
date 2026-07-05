package repo

import (
	"context"
	"errors"

	model "backend/model/schema"
	"backend/model/view"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ===================== OrderRepo =====================

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

// Create 创建订单及入住人（事务）。
func (r *OrderRepo) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Guests").Create(order).Error; err != nil {
			return err
		}
		if len(order.Guests) > 0 {
			if err := tx.CreateInBatches(order.Guests, len(order.Guests)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// FindByID 根据订单 ID 查询，预加载 User、Room、Guests（及 Guests.Person）。
func (r *OrderRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Room").
		Preload("Guests").
		Preload("Guests.Person").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// FindByUserID 根据用户 ID 查询订单列表，预加载 Room，按 create_at 降序。
func (r *OrderRepo) FindByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]model.Order, int64, error) {
	var results []model.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Preload("Room").Order("create_at DESC").Find(&results).Error
	return results, total, err
}

// FindByUserIDAndStatus 根据用户 ID 和状态查询订单列表，预加载 Room，按 create_at 降序。
func (r *OrderRepo) FindByUserIDAndStatus(ctx context.Context, userID uuid.UUID, status model.OrderStatus, offset, limit int) ([]model.Order, int64, error) {
	var results []model.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ? AND status = ?", userID, status)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Preload("Room").Order("create_at DESC").Find(&results).Error
	return results, total, err
}

// FindByHotelID 根据酒店 ID 查询订单列表，通过 room_1718 关联，预加载 User，按 create_at 降序。
func (r *OrderRepo) FindByHotelID(ctx context.Context, hotelID uuid.UUID, offset, limit int) ([]model.Order, int64, error) {
	var results []model.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Order{}).
		Joins("JOIN room_1718 ON room_1718.id = order_1718.room_id").
		Where("room_1718.hotel_id = ?", hotelID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Preload("User").Order("create_at DESC").Find(&results).Error
	return results, total, err
}

// FindAll 查询全部订单，预加载 User、Room，按 create_at 降序。
func (r *OrderRepo) FindAll(ctx context.Context, offset, limit int) ([]model.Order, int64, error) {
	var results []model.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Order{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Preload("User").Preload("Room").Order("create_at DESC").Find(&results).Error
	return results, total, err
}

// Update 更新订单全部字段。
func (r *OrderRepo) Update(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// UpdateStatus 更新订单状态，禁止非法状态流转。
func (r *OrderRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) error {
	var order model.Order
	if err := r.db.WithContext(ctx).First(&order, "id = ?", id).Error; err != nil {
		return err
	}

	invalidTransitions := map[model.OrderStatus]model.OrderStatus{
		model.OrderCancelled: model.OrderCompleted,
		model.OrderCompleted: model.OrderPending,
		model.OrderBooked:    model.OrderPending,
	}

	if target, ok := invalidTransitions[order.Status]; ok && target == status {
		return errors.New("invalid status transition")
	}

	return r.db.WithContext(ctx).Model(&order).Update("status", status).Error
}

// Delete 根据 ID 硬删除订单。
func (r *OrderRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Order{}, "id = ?", id).Error
}

// ===================== OrderGuestRepo =====================

type OrderGuestRepo struct {
	db *gorm.DB
}

func NewOrderGuestRepo(db *gorm.DB) *OrderGuestRepo {
	return &OrderGuestRepo{db: db}
}

// Create 插入单条入住人记录。
func (r *OrderGuestRepo) Create(ctx context.Context, guest *model.OrderGuest) error {
	return r.db.WithContext(ctx).Create(guest).Error
}

// BatchCreate 批量插入入住人记录。
func (r *OrderGuestRepo) BatchCreate(ctx context.Context, guests []model.OrderGuest) error {
	return r.db.WithContext(ctx).CreateInBatches(guests, len(guests)).Error
}

// FindByOrderID 根据订单 ID 查询入住人列表，预加载 Person。
func (r *OrderGuestRepo) FindByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.OrderGuest, error) {
	var guests []model.OrderGuest
	err := r.db.WithContext(ctx).Preload("Person").Where("order_id = ?", orderID).Find(&guests).Error
	return guests, err
}

// FindByIDCard 根据身份证号查询入住人列表，预加载 Order，支持分页。
func (r *OrderGuestRepo) FindByIDCard(ctx context.Context, idCard string, offset, limit int) ([]model.OrderGuest, int64, error) {
	var results []model.OrderGuest
	var total int64
	query := r.db.WithContext(ctx).Model(&model.OrderGuest{}).Where("id_card = ?", idCard)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Preload("Order").Find(&results).Error
	return results, total, err
}

// Delete 根据订单 ID 和身份证号删除入住人记录。
func (r *OrderGuestRepo) Delete(ctx context.Context, orderID uuid.UUID, idCard string) error {
	return r.db.WithContext(ctx).
		Where("order_id = ? AND id_card = ?", orderID, idCard).
		Delete(&model.OrderGuest{}).Error
}

// ===================== OrderFullRepo =====================

// OrderFullRepo 订单完整视图（view_order_full_1718）只读仓库。
type OrderFullRepo struct {
	db *gorm.DB
}

func NewOrderFullRepo(db *gorm.DB) *OrderFullRepo {
	return &OrderFullRepo{db: db}
}

// FindByOrderID 根据订单 ID 查询视图记录，按 guest_id_card 排序。
func (r *OrderFullRepo) FindByOrderID(ctx context.Context, orderID uuid.UUID) ([]view.OrderFull, error) {
	var results []view.OrderFull
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("guest_id_card").
		Find(&results).Error
	return results, err
}

// FindByUserID 根据用户 ID 查询视图记录，按 create_at 降序，支持分页。
func (r *OrderFullRepo) FindByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]view.OrderFull, int64, error) {
	var results []view.OrderFull
	var total int64
	query := r.db.WithContext(ctx).Model(&view.OrderFull{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Order("create_at DESC").Find(&results).Error
	return results, total, err
}

// FindByHotelID 根据酒店 ID 查询视图记录，按 create_at 降序，支持分页。
func (r *OrderFullRepo) FindByHotelID(ctx context.Context, hotelID uuid.UUID, offset, limit int) ([]view.OrderFull, int64, error) {
	var results []view.OrderFull
	var total int64
	query := r.db.WithContext(ctx).Model(&view.OrderFull{}).Where("hotel_id = ?", hotelID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Order("create_at DESC").Find(&results).Error
	return results, total, err
}

// ===================== MyOrdersRepo =====================

// MyOrdersRepo 我的订单视图（view_my_orders_1718）只读仓库。
type MyOrdersRepo struct {
	db *gorm.DB
}

func NewMyOrdersRepo(db *gorm.DB) *MyOrdersRepo {
	return &MyOrdersRepo{db: db}
}

// FindByUserID 根据用户 ID 查询我的订单列表，按 create_at 降序，支持分页。
func (r *MyOrdersRepo) FindByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]view.MyOrders, int64, error) {
	var results []view.MyOrders
	var total int64
	query := r.db.WithContext(ctx).Model(&view.MyOrders{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Order("create_at DESC").Find(&results).Error
	return results, total, err
}

// FindByUserIDAndStatus 根据用户 ID 和状态查询我的订单列表，按 create_at 降序，支持分页。
func (r *MyOrdersRepo) FindByUserIDAndStatus(ctx context.Context, userID uuid.UUID, status string, offset, limit int) ([]view.MyOrders, int64, error) {
	var results []view.MyOrders
	var total int64
	query := r.db.WithContext(ctx).Model(&view.MyOrders{}).Where("user_id = ? AND order_status = ?", userID, status)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Order("create_at DESC").Find(&results).Error
	return results, total, err
}

// FindByOrderID 根据订单 ID 查询单条我的订单记录。
func (r *MyOrdersRepo) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*view.MyOrders, error) {
	var result view.MyOrders
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
