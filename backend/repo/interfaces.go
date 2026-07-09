// Package repo 定义数据访问层接口，用于依赖注入与单元测试 mock。
package repo

import (
	"context"

	model "backend/model/schema"
	"backend/model/view"

	"github.com/google/uuid"
)

// ─── 核心实体 Repository 接口 ──────────────────────────────────

// UserRepository 用户数据访问接口。
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	FindAll(ctx context.Context, offset, limit int, role *model.UserRole) ([]model.User, int64, error)
	Update(ctx context.Context, user *model.User) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error
	UpdatePoints(ctx context.Context, userID uuid.UUID, points int32) error
	UpdateVipLevel(ctx context.Context, userID uuid.UUID, vipLevel int16) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// VipLevelRepository VIP 等级数据访问接口。
type VipLevelRepository interface {
	FindByLevel(ctx context.Context, level int16) (*model.VipLevel, error)
	FindAll(ctx context.Context) ([]model.VipLevel, error)
	FindByMinPoints(ctx context.Context, points int32) (*model.VipLevel, error)
	Create(ctx context.Context, vip *model.VipLevel) error
	Update(ctx context.Context, vip *model.VipLevel) error
	Delete(ctx context.Context, level int16) error
}

// HotelRepository 酒店数据访问接口。
type HotelRepository interface {
	Create(ctx context.Context, hotel *model.Hotel) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Hotel, error)
	FindAll(ctx context.Context, offset, limit int, regionID *int, starLevel *int16, keyword string) ([]model.Hotel, int64, error)
	Update(ctx context.Context, hotel *model.Hotel) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// RoomRepository 客房数据访问接口。
type RoomRepository interface {
	Create(ctx context.Context, room *model.Room) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Room, error)
	FindByHotelID(ctx context.Context, hotelID uuid.UUID, offset, limit int) ([]model.Room, int64, error)
	FindAll(ctx context.Context, offset, limit int) ([]model.Room, int64, error)
	Update(ctx context.Context, room *model.Room) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// OrderRepository 订单数据访问接口。
type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]model.Order, int64, error)
	FindByUserIDAndStatus(ctx context.Context, userID uuid.UUID, status model.OrderStatus, offset, limit int) ([]model.Order, int64, error)
	FindByHotelID(ctx context.Context, hotelID uuid.UUID, offset, limit int) ([]model.Order, int64, error)
	FindAll(ctx context.Context, offset, limit int) ([]model.Order, int64, error)
	Update(ctx context.Context, order *model.Order) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ReviewRepository 评价数据访问接口。
type ReviewRepository interface {
	Create(ctx context.Context, review *model.Review) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Review, error)
	FindByHotelID(ctx context.Context, hotelID uuid.UUID, offset, limit int) ([]model.Review, int64, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]model.Review, int64, error)
	FindByOrderID(ctx context.Context, orderID uuid.UUID) (*model.Review, error)
	FindAll(ctx context.Context, offset, limit int) ([]model.Review, int64, error)
	FindByRating(ctx context.Context, rating int16, offset, limit int) ([]model.Review, int64, error)
	Update(ctx context.Context, review *model.Review) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// PersonRepository 入住人员数据访问接口。
type PersonRepository interface {
	Create(ctx context.Context, person *model.Person) error
	FindByIDCard(ctx context.Context, idCard string) (*model.Person, error)
	FindAll(ctx context.Context, offset, limit int, keyword string) ([]model.Person, int64, error)
	Update(ctx context.Context, person *model.Person) error
	Delete(ctx context.Context, idCard string) error
	Upsert(ctx context.Context, person *model.Person) error
}

// RegionRepository 地区数据访问接口。
type RegionRepository interface {
	FindByID(ctx context.Context, id int) (*model.Region, error)
	FindByParentID(ctx context.Context, parentID int) ([]model.Region, error)
	FindAllProvinces(ctx context.Context) ([]model.Region, error)
	FindAll(ctx context.Context) ([]model.Region, error)
	Create(ctx context.Context, region *model.Region) error
	Update(ctx context.Context, region *model.Region) error
	Delete(ctx context.Context, id int) error
}

// ─── 视图 Repository 接口（只读）────────────────────────────────

// HotelSummaryRepository 酒店摘要视图接口。
type HotelSummaryRepository interface {
	FindByID(ctx context.Context, hotelID uuid.UUID) (*view.HotelSummary, error)
	FindAll(ctx context.Context, offset, limit int, province, city, district string, starLevel *int16, minPrice, maxPrice *float64) ([]view.HotelSummary, int64, error)
	FindByRegionID(ctx context.Context, regionID int, offset, limit int) ([]view.HotelSummary, int64, error)
}

// RoomDetailsRepository 客房详情视图接口。
type RoomDetailsRepository interface {
	FindByRoomID(ctx context.Context, roomID uuid.UUID) (*view.RoomDetails, error)
	FindByHotelID(ctx context.Context, hotelID uuid.UUID, offset, limit int) ([]view.RoomDetails, int64, error)
	FindAll(ctx context.Context, offset, limit int, province, city, district string, starLevel *int16, minPrice, maxPrice *float64) ([]view.RoomDetails, int64, error)
}

// OrderFullRepository 订单完整视图接口。
type OrderFullRepository interface {
	FindByOrderID(ctx context.Context, orderID uuid.UUID) ([]view.OrderFull, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]view.OrderFull, int64, error)
	FindByHotelID(ctx context.Context, hotelID uuid.UUID, offset, limit int) ([]view.OrderFull, int64, error)
}

// ReviewFullRepository 评价完整视图接口。
type ReviewFullRepository interface {
	FindByReviewID(ctx context.Context, reviewID uuid.UUID) (*view.ReviewFull, error)
	FindByHotelID(ctx context.Context, hotelID uuid.UUID, offset, limit int) ([]view.ReviewFull, int64, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]view.ReviewFull, int64, error)
	FindAll(ctx context.Context, offset, limit int) ([]view.ReviewFull, int64, error)
}

// UserVipRepository 用户 VIP 视图接口。
type UserVipRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*view.UserVip, error)
	FindAll(ctx context.Context, offset, limit int, role string) ([]view.UserVip, int64, error)
	FindByVipLevel(ctx context.Context, level int16, offset, limit int) ([]view.UserVip, int64, error)
}

// PersonInfoRepository 入住人信息视图接口。
type PersonInfoRepository interface {
	FindByIDCard(ctx context.Context, idCard string) (*view.PersonInfo, error)
	FindAll(ctx context.Context, offset, limit int, gender string, minAge, maxAge *int) ([]view.PersonInfo, int64, error)
}

// GuestBookingStatsRepository 入住人统计视图接口。
type GuestBookingStatsRepository interface {
	FindByIDCard(ctx context.Context, idCard string) (*view.GuestBookingStats, error)
	FindAll(ctx context.Context, offset, limit int, ageGroup, gender, favCity string) ([]view.GuestBookingStats, int64, error)
	FindTopGuests(ctx context.Context, limit int) ([]view.GuestBookingStats, error)
}

// MyOrdersRepository 我的订单视图接口。
type MyOrdersRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]view.MyOrders, int64, error)
	FindByUserIDAndStatus(ctx context.Context, userID uuid.UUID, status string, offset, limit int) ([]view.MyOrders, int64, error)
	FindByOrderID(ctx context.Context, orderID uuid.UUID) (*view.MyOrders, error)
}

// OrderDetailRepository 订单详情视图接口。
type OrderDetailRepository interface {
	FindDetailByOrderID(ctx context.Context, orderID uuid.UUID) (*view.OrderDetail, error)
}

// OrderSummaryRepository 订单概览视图接口。
type OrderSummaryRepository interface {
	FindAll(ctx context.Context, offset, limit int) ([]view.OrderSummary, int64, error)
	FindByStatus(ctx context.Context, status string, offset, limit int) ([]view.OrderSummary, int64, error)
}

// ─── 编译时类型检查 ──────────────────────────────────────────
// 确保所有具体 repo 结构体满足对应接口（IDE 提示 + 编译器保证）

var (
	_ UserRepository              = (*UserRepo)(nil)
	_ VipLevelRepository          = (*VipLevelRepo)(nil)
	_ HotelRepository             = (*HotelRepo)(nil)
	_ RoomRepository              = (*RoomRepo)(nil)
	_ OrderRepository             = (*OrderRepo)(nil)
	_ ReviewRepository            = (*ReviewRepo)(nil)
	_ PersonRepository            = (*PersonRepo)(nil)
	_ RegionRepository            = (*RegionRepo)(nil)
	_ HotelSummaryRepository      = (*HotelSummaryRepo)(nil)
	_ RoomDetailsRepository       = (*RoomDetailsRepo)(nil)
	_ OrderFullRepository         = (*OrderFullRepo)(nil)
	_ ReviewFullRepository        = (*ReviewFullRepo)(nil)
	_ UserVipRepository           = (*UserVipRepo)(nil)
	_ PersonInfoRepository        = (*PersonInfoRepo)(nil)
	_ GuestBookingStatsRepository = (*GuestBookingStatsRepo)(nil)
	_ MyOrdersRepository          = (*MyOrdersRepo)(nil)
	_ OrderDetailRepository       = (*OrderDetailRepo)(nil)
	_ OrderSummaryRepository      = (*OrderSummaryRepo)(nil)
)
