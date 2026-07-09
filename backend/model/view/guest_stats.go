package view

import (
	"time"
)

// GuestBookingStats 对应视图 view_guest_booking_stats_1718（只读）
//
// 入住人预订统计分析，按年龄/性别/职业/学历/收入等维度分析客户偏好。
// total_orders 只统计 booked/checked_in/completed 状态的订单，
// 与 total_amount、avg_order_amount 口径一致。
type GuestBookingStats struct {
	PersonIDCard   string     `json:"personIdCard"   gorm:"column:person_id_card;type:char(18);->"`
	PersonName     string     `json:"personName"     gorm:"column:person_name;type:text;->"`
	Gender         *string    `json:"gender"         gorm:"column:gender;type:text;->"`
	Age            *int       `json:"age"            gorm:"column:age;type:integer;->"`
	AgeGroup       string     `json:"ageGroup"       gorm:"column:age_group;type:text;->"`
	TotalOrders    int64      `json:"totalOrders"    gorm:"column:total_orders;type:bigint;->"`
	TotalNights    int64      `json:"totalNights"    gorm:"column:total_nights;type:bigint;->"`
	TotalAmount    float64    `json:"totalAmount"    gorm:"column:total_amount;type:numeric;->"`
	AvgOrderAmount float64    `json:"avgOrderAmount" gorm:"column:avg_order_amount;type:numeric;->"`
	FavCity        *string    `json:"favCity"        gorm:"column:fav_city;type:text;->"`
	FavHotel       *string    `json:"favHotel"       gorm:"column:fav_hotel;type:text;->"`
	FavRoomType    *string    `json:"favRoomType"    gorm:"column:fav_room_type;type:text;->"`
	LastOrderDate  *time.Time `json:"lastOrderDate"  gorm:"column:last_order_date;type:timestamptz;->"`
}

func (GuestBookingStats) TableName() string {
	return "view_guest_booking_stats_1718"
}
