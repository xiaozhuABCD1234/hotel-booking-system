package view

import (
	"github.com/google/uuid"
)

// RoomDetails 对应视图 view_room_details_1718（只读）
//
// 按城市 → 区域 → 酒店查看所有客房详细信息。
// 所有字段标记 gorm:"->" 表示只读，GORM 不会将其纳入写入操作。
type RoomDetails struct {
	Province          *string   `gorm:"column:province;type:text;->"`
	City              *string   `gorm:"column:city;type:text;->"`
	District          string    `gorm:"column:district;type:text;->"`
	HotelID           uuid.UUID `gorm:"column:hotel_id;type:uuid;->"`
	HotelName         string    `gorm:"column:hotel_name;type:text;->"`
	Address           string    `gorm:"column:address;type:text;->"`
	StarLevel         *int16    `gorm:"column:star_level;type:smallint;->"`
	HotelTelephone    string    `gorm:"column:hotel_telephone;type:varchar(20);->"`
	RoomID            uuid.UUID `gorm:"column:room_id;type:uuid;->"`
	TypeName          string    `gorm:"column:type_name;type:text;->"`
	TotalQuantity     int32     `gorm:"column:total_quantity;type:integer;->"`
	AvailableQuantity int32     `gorm:"column:available_quantity;type:integer;->"`
	Price             float64   `gorm:"column:price;type:numeric(10,2);->"`
	WeekendPrice      *float64  `gorm:"column:weekend_price;type:numeric(10,2);->"`
	RoomDescription   *string   `gorm:"column:room_description;type:text;->"`
	Facilities        *string   `gorm:"column:facilities;type:text;->"`
	AvgRating         *float64  `gorm:"column:avg_rating;type:numeric(3,2);->"`
	ReviewCount       *int64    `gorm:"column:review_count;type:bigint;->"`
}

func (RoomDetails) TableName() string {
	return "view_room_details_1718"
}
