package view

import (
	"time"

	"github.com/google/uuid"
)

// MyOrders 对应视图 view_my_orders_1718（只读）
//
// 我的订单列表，每个订单一行，不展开入住人。
// OrderStatus 以 string 接收 order_status 枚举值。
type MyOrders struct {
	OrderID      uuid.UUID `json:"orderId"      gorm:"column:order_id;type:uuid;->"`
	UserID       uuid.UUID `json:"userId"       gorm:"column:user_id;type:uuid;->"`
	HotelName    string    `json:"hotelName"    gorm:"column:hotel_name;type:text;->"`
	City         *string   `json:"city"         gorm:"column:city;type:text;->"`
	RoomType     string    `json:"roomType"     gorm:"column:room_type;type:text;->"`
	Quantity     int32     `json:"quantity"     gorm:"column:quantity;type:integer;->"`
	CheckInDate  time.Time `json:"checkInDate"  gorm:"column:check_in_date;type:date;->"`
	CheckOutDate time.Time `json:"checkOutDate" gorm:"column:check_out_date;type:date;->"`
	Nights       int       `json:"nights"       gorm:"column:nights;type:integer;->"`
	ActualPrice  float64   `json:"actualPrice"  gorm:"column:actual_price;type:numeric(10,2);->"`
	OrderStatus  string    `json:"orderStatus"  gorm:"column:order_status;type:order_status;->"`
	GuestCount   int64     `json:"guestCount"   gorm:"column:guest_count;type:bigint;->"`
	CreateAt     time.Time `json:"createAt"     gorm:"column:create_at;type:timestamptz;->"`
}

func (MyOrders) TableName() string {
	return "view_my_orders_1718"
}
