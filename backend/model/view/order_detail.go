package view

import (
	"time"

	"github.com/google/uuid"
)

// OrderDetail 对应视图 view_order_detail_1718（只读）
//
// 订单详情，下单人与入住人明确区分，入住人聚合为一个字段。
type OrderDetail struct {
	OrderID        uuid.UUID `json:"orderId"        gorm:"column:order_id;type:uuid;->"`
	Status         string    `json:"status"         gorm:"column:status;type:varchar(20);->"`
	Quantity       int32     `json:"quantity"       gorm:"column:quantity;type:integer;->"`
	CheckInDate    time.Time `json:"checkInDate"    gorm:"column:check_in_date;type:date;->"`
	CheckOutDate   time.Time `json:"checkOutDate"   gorm:"column:check_out_date;type:date;->"`
	Nights         int       `json:"nights"         gorm:"column:nights;type:integer;->"`
	TotalPrice     float64   `json:"totalPrice"     gorm:"column:total_price;type:numeric(10,2);->"`
	ActualPrice    float64   `json:"actualPrice"    gorm:"column:actual_price;type:numeric(10,2);->"`
	CreateAt       time.Time `json:"createAt"       gorm:"column:create_at;type:timestamptz;->"`
	OrderUser      string    `json:"orderUser"      gorm:"column:order_user;type:text;->"`
	OrderUserName  *string   `json:"orderUserName"  gorm:"column:order_user_name;type:text;->"`
	OrderUserPhone *string   `json:"orderUserPhone" gorm:"column:order_user_phone;type:varchar(20);->"`
	HotelName      string    `json:"hotelName"      gorm:"column:hotel_name;type:text;->"`
	City           *string   `json:"city"           gorm:"column:city;type:text;->"`
	RoomType       string    `json:"roomType"       gorm:"column:room_type;type:text;->"`
	RoomPrice      float64   `json:"roomPrice"      gorm:"column:room_price;type:numeric(10,2);->"`
	GuestCount     int64     `json:"guestCount"     gorm:"column:guest_count;type:bigint;->"`
	GuestNames     string    `json:"guestNames"     gorm:"column:guest_names;type:text;->"`
}

func (OrderDetail) TableName() string {
	return "view_order_detail_1718"
}
