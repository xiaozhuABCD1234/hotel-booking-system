package view

import (
	"github.com/google/uuid"
)

// RoomDetails 对应视图 view_room_details_1718（只读）
//
// 按城市 → 区域 → 酒店查看所有客房详细信息。
// 所有字段标记 gorm:"->" 表示只读，GORM 不会将其纳入写入操作。
type RoomDetails struct {
	Province          *string   `json:"province"          gorm:"column:province;type:text;->"`
	City              *string   `json:"city"              gorm:"column:city;type:text;->"`
	District          string    `json:"district"          gorm:"column:district;type:text;->"`
	HotelID           uuid.UUID `json:"hotelId"           gorm:"column:hotel_id;type:uuid;->"`
	HotelName         string    `json:"hotelName"         gorm:"column:hotel_name;type:text;->"`
	Address           string    `json:"address"           gorm:"column:address;type:text;->"`
	StarLevel         *int16    `json:"starLevel"         gorm:"column:star_level;type:smallint;->"`
	HotelTelephone    string    `json:"hotelTelephone"    gorm:"column:hotel_telephone;type:varchar(20);->"`
	RoomID            uuid.UUID `json:"roomId"            gorm:"column:room_id;type:uuid;->"`
	TypeName          string    `json:"typeName"          gorm:"column:type_name;type:text;->"`
	TotalQuantity     int32     `json:"totalQuantity"     gorm:"column:total_quantity;type:integer;->"`
	AvailableQuantity int32     `json:"availableQuantity" gorm:"column:available_quantity;type:integer;->"`
	Price             float64   `json:"price"             gorm:"column:price;type:numeric(10,2);->"`
	WeekendPrice      *float64  `json:"weekendPrice"      gorm:"column:weekend_price;type:numeric(10,2);->"`
	RoomDescription   *string   `json:"roomDescription"   gorm:"column:room_description;type:text;->"`
	Facilities        *string   `json:"facilities"        gorm:"column:facilities;type:text;->"`
	AvgRating         *float64  `json:"avgRating"         gorm:"column:avg_rating;type:numeric(3,2);->"`
	ReviewCount       *int64    `json:"reviewCount"       gorm:"column:review_count;type:bigint;->"`
}

func (RoomDetails) TableName() string {
	return "view_room_details_1718"
}
