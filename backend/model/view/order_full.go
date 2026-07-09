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
	OrderID         uuid.UUID `json:"orderId"         gorm:"column:order_id;type:uuid;->"`
	UserID          uuid.UUID `json:"userId"          gorm:"column:user_id;type:uuid;->"`
	Username        string    `json:"username"        gorm:"column:username;type:text;->"`
	UserPhone       *string   `json:"userPhone"       gorm:"column:user_phone;type:varchar(20);->"`
	UserRealName    *string   `json:"userRealName"    gorm:"column:user_real_name;type:text;->"`
	HotelID         uuid.UUID `json:"hotelId"         gorm:"column:hotel_id;type:uuid;->"`
	HotelName       string    `json:"hotelName"       gorm:"column:hotel_name;type:text;->"`
	City            *string   `json:"city"            gorm:"column:city;type:text;->"`
	District        string    `json:"district"        gorm:"column:district;type:text;->"`
	HotelTelephone  string    `json:"hotelTelephone"  gorm:"column:hotel_telephone;type:varchar(20);->"`
	RoomID          uuid.UUID `json:"roomId"          gorm:"column:room_id;type:uuid;->"`
	RoomType        string    `json:"roomType"        gorm:"column:room_type;type:text;->"`
	Quantity        int32     `json:"quantity"        gorm:"column:quantity;type:integer;->"`
	CheckInDate     time.Time `json:"checkInDate"     gorm:"column:check_in_date;type:date;->"`
	CheckOutDate    time.Time `json:"checkOutDate"    gorm:"column:check_out_date;type:date;->"`
	Nights          int       `json:"nights"          gorm:"column:nights;type:integer;->"`
	TotalPrice      float64   `json:"totalPrice"      gorm:"column:total_price;type:numeric(10,2);->"`
	ActualPrice     float64   `json:"actualPrice"     gorm:"column:actual_price;type:numeric(10,2);->"`
	VipDiscountRate float64   `json:"vipDiscountRate" gorm:"column:vip_discount_rate;type:numeric(3,2);->"`
	OrderStatus     string    `json:"orderStatus"     gorm:"column:order_status;type:order_status;->"`
	GuestIDCard     *string   `json:"guestIdCard"     gorm:"column:guest_id_card;type:char(18);->"`
	GuestName       *string   `json:"guestName"       gorm:"column:guest_name;type:text;->"`
	GuestGender     *string   `json:"guestGender"     gorm:"column:guest_gender;type:text;->"`
	GuestAge        *int      `json:"guestAge"        gorm:"column:guest_age;type:integer;->"`
	CreateAt        time.Time `json:"createAt"        gorm:"column:create_at;type:timestamptz;->"`
}

func (OrderFull) TableName() string {
	return "view_order_full_1718"
}
