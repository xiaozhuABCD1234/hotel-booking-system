package view

import (
	"time"

	"github.com/google/uuid"
)

// OrderDetail 对应视图 view_order_detail_1718（只读）
//
// 订单详情，下单人与入住人明确区分，入住人聚合为一个字段。
type OrderDetail struct {
	OrderID        uuid.UUID `gorm:"column:order_id;type:uuid;->"`
	Status         string    `gorm:"column:status;type:varchar(20);->"`
	Quantity       int32     `gorm:"column:quantity;type:integer;->"`
	CheckInDate    time.Time `gorm:"column:check_in_date;type:date;->"`
	CheckOutDate   time.Time `gorm:"column:check_out_date;type:date;->"`
	Nights         int       `gorm:"column:nights;type:integer;->"`
	TotalPrice     float64   `gorm:"column:total_price;type:numeric(10,2);->"`
	ActualPrice    float64   `gorm:"column:actual_price;type:numeric(10,2);->"`
	CreateAt       time.Time `gorm:"column:create_at;type:timestamptz;->"`
	OrderUser      string    `gorm:"column:order_user;type:text;->"`
	OrderUserName  *string   `gorm:"column:order_user_name;type:text;->"`
	OrderUserPhone *string   `gorm:"column:order_user_phone;type:varchar(20);->"`
	HotelName      string    `gorm:"column:hotel_name;type:text;->"`
	City           *string   `gorm:"column:city;type:text;->"`
	RoomType       string    `gorm:"column:room_type;type:text;->"`
	RoomPrice      float64   `gorm:"column:room_price;type:numeric(10,2);->"`
	GuestCount     int64     `gorm:"column:guest_count;type:bigint;->"`
	GuestNames     string    `gorm:"column:guest_names;type:text;->"`
}

func (OrderDetail) TableName() string {
	return "view_order_detail_1718"
}
