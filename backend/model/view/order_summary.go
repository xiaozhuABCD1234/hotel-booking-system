package view

import (
	"time"

	"github.com/google/uuid"
)

// OrderSummary 对应视图 view_order_summary_1718（只读）
//
// 订单概览，字段精简，适合管理端列表页。
type OrderSummary struct {
	OrderID       uuid.UUID `json:"orderId"       gorm:"column:order_id;type:uuid;->"`
	Status        string    `json:"status"        gorm:"column:status;type:varchar(20);->"`
	Quantity      int32     `json:"quantity"      gorm:"column:quantity;type:integer;->"`
	CheckInDate   time.Time `json:"checkInDate"   gorm:"column:check_in_date;type:date;->"`
	CheckOutDate  time.Time `json:"checkOutDate"  gorm:"column:check_out_date;type:date;->"`
	Nights        int       `json:"nights"        gorm:"column:nights;type:integer;->"`
	ActualPrice   float64   `json:"actualPrice"   gorm:"column:actual_price;type:numeric(10,2);->"`
	CreateAt      time.Time `json:"createAt"      gorm:"column:create_at;type:timestamptz;->"`
	OrderUserName *string   `json:"orderUserName" gorm:"column:order_user_name;type:text;->"`
	HotelName     string    `json:"hotelName"     gorm:"column:hotel_name;type:text;->"`
	RoomType      string    `json:"roomType"      gorm:"column:room_type;type:text;->"`
	GuestCount    int64     `json:"guestCount"    gorm:"column:guest_count;type:bigint;->"`
}

func (OrderSummary) TableName() string {
	return "view_order_summary_1718"
}
