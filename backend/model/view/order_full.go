package view

import (
	"time"

	"github.com/google/uuid"
)

// OrderFull 对应视图 view_order_full_1718（只读）
//
// 订单完整详情，一单多个入住人时每个入住人一行。
// OrderStatus 以 string 接收 order_status 枚举值。
type OrderFull struct {
	OrderID         uuid.UUID `gorm:"column:order_id;type:uuid;->"`
	UserID          uuid.UUID `gorm:"column:user_id;type:uuid;->"`
	Username        string    `gorm:"column:username;type:text;->"`
	UserPhone       *string   `gorm:"column:user_phone;type:varchar(20);->"`
	UserRealName    *string   `gorm:"column:user_real_name;type:text;->"`
	HotelID         uuid.UUID `gorm:"column:hotel_id;type:uuid;->"`
	HotelName       string    `gorm:"column:hotel_name;type:text;->"`
	City            *string   `gorm:"column:city;type:text;->"`
	District        string    `gorm:"column:district;type:text;->"`
	HotelTelephone  string    `gorm:"column:hotel_telephone;type:varchar(20);->"`
	RoomID          uuid.UUID `gorm:"column:room_id;type:uuid;->"`
	RoomType        string    `gorm:"column:room_type;type:text;->"`
	Quantity        int32     `gorm:"column:quantity;type:integer;->"`
	CheckInDate     time.Time `gorm:"column:check_in_date;type:date;->"`
	CheckOutDate    time.Time `gorm:"column:check_out_date;type:date;->"`
	Nights          int       `gorm:"column:nights;type:integer;->"`
	TotalPrice      float64   `gorm:"column:total_price;type:numeric(10,2);->"`
	ActualPrice     float64   `gorm:"column:actual_price;type:numeric(10,2);->"`
	VipDiscountRate float64   `gorm:"column:vip_discount_rate;type:numeric(3,2);->"`
	OrderStatus     string    `gorm:"column:order_status;type:order_status;->"`
	GuestIDCard     *string   `gorm:"column:guest_id_card;type:char(18);->"`
	GuestName       *string   `gorm:"column:guest_name;type:text;->"`
	GuestGender     *string   `gorm:"column:guest_gender;type:text;->"`
	GuestAge        *int      `gorm:"column:guest_age;type:integer;->"`
	CreateAt        time.Time `gorm:"column:create_at;type:timestamptz;->"`
}

func (OrderFull) TableName() string {
	return "view_order_full_1718"
}
