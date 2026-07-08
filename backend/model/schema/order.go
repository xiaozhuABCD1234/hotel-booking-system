package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderStatus 订单状态枚举，对应 PostgreSQL 自定义类型 order_status。
type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderBooked    OrderStatus = "booked"
	OrderCheckedIn OrderStatus = "checked_in"
	OrderCompleted OrderStatus = "completed"
	OrderCancelled OrderStatus = "cancelled"
)

// Order 订单表，对应表 order_1718
//
// 注意：表级 CHECK (check_out_date > check_in_date) 在 BeforeSave 中校验。
type Order struct {
	ID           uuid.UUID    `json:"id"            gorm:"column:id;type:uuid;primaryKey;default:uuidv7()"`
	UserID       uuid.UUID    `json:"userId"        gorm:"column:user_id;type:uuid;not null;index:idx_order_user"`
	User         User         `json:"user"          gorm:"foreignKey:UserID;references:ID"`
	RoomID       uuid.UUID    `json:"roomId"        gorm:"column:room_id;type:uuid;not null;index:idx_order_room"`
	Room         Room         `json:"room"          gorm:"foreignKey:RoomID;references:ID"`
	Quantity     int32        `json:"quantity"      gorm:"column:quantity;type:integer;not null;check:quantity > 0"`
	CheckInDate  time.Time    `json:"checkInDate"   gorm:"column:check_in_date;type:date;not null;index:idx_order_dates,priority:1"`
	CheckOutDate time.Time    `json:"checkOutDate"  gorm:"column:check_out_date;type:date;not null;index:idx_order_dates,priority:2"`
	TotalPrice   float64      `json:"totalPrice"    gorm:"column:total_price;type:numeric(10,2);not null;check:total_price >= 0"`
	Discount     float64      `json:"discount"      gorm:"column:discount;type:numeric(10,2);not null;default:0;check:discount >= 0"`
	ActualPrice  float64      `json:"actualPrice"   gorm:"column:actual_price;type:numeric(10,2);not null;check:actual_price >= 0"`
	Status       OrderStatus  `json:"status"        gorm:"column:status;type:order_status;not null;default:pending;index:idx_order_status"`
	CreateAt     time.Time    `json:"createAt"      gorm:"column:create_at;type:timestamptz;autoCreateTime;index:idx_order_create"`
	UpdateAt     time.Time    `json:"updateAt"      gorm:"column:update_at;type:timestamptz;autoUpdateTime"`
	Guests       []OrderGuest `json:"guests"        gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Order) TableName() string {
	return "order_1718"
}

func (o *Order) BeforeSave(tx *gorm.DB) error {
	if !o.CheckOutDate.After(o.CheckInDate) {
		return errors.New("check_out_date must be after check_in_date")
	}
	return nil
}

// OrderGuest 入住人员关联表（订单 ↔ 入住人，多对多），对应表 order_guest_1718
type OrderGuest struct {
	OrderID uuid.UUID `json:"orderId" gorm:"column:order_id;type:uuid;primaryKey"`
	IDCard  string    `json:"idCard"  gorm:"column:id_card;type:varchar(18);primaryKey;index:idx_order_guest_id_card"`
	Order   Order     `json:"-"       gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE"`
	Person  Person    `json:"person"  gorm:"foreignKey:IDCard;references:IDCard"`
}

func (OrderGuest) TableName() string {
	return "order_guest_1718"
}
