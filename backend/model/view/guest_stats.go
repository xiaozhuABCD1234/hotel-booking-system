package view

import (
	"time"
)

// GuestBookingStats 对应视图 view_guest_booking_stats_1718（只读）
//
// 入住人预订统计分析，按年龄/性别等维度分析客户偏好。
// total_orders 只统计 booked/checked_in/completed 状态的订单，
// 与 total_amount、avg_order_amount 口径一致。
type GuestBookingStats struct {
	PersonIDCard   string     `gorm:"column:person_id_card;type:varchar(18);->"`
	PersonName     string     `gorm:"column:person_name;type:text;->"`
	Gender         *string    `gorm:"column:gender;type:text;->"`
	Age            *int       `gorm:"column:age;type:integer;->"`
	AgeGroup       string     `gorm:"column:age_group;type:text;->"`
	TotalOrders    int64      `gorm:"column:total_orders;type:bigint;->"`
	TotalNights    int64      `gorm:"column:total_nights;type:bigint;->"`
	TotalAmount    float64    `gorm:"column:total_amount;type:numeric;->"`
	AvgOrderAmount float64    `gorm:"column:avg_order_amount;type:numeric;->"`
	FavCity        *string    `gorm:"column:fav_city;type:text;->"`
	FavHotel       *string    `gorm:"column:fav_hotel;type:text;->"`
	FavRoomType    *string    `gorm:"column:fav_room_type;type:text;->"`
	LastOrderDate  *time.Time `gorm:"column:last_order_date;type:timestamptz;->"`
}

func (GuestBookingStats) TableName() string {
	return "view_guest_booking_stats_1718"
}
