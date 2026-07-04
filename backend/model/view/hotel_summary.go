package view

import (
	"github.com/google/uuid"
)

// HotelSummary 对应视图 view_hotel_summary_1718（只读）
//
// 酒店摘要列表，用于搜索页、首页推荐。
// min_price / avg_rating 经 COALESCE 兜底为 0，不会为 NULL。
type HotelSummary struct {
	HotelID     uuid.UUID `gorm:"column:hotel_id;type:uuid;->"`
	HotelName   string    `gorm:"column:hotel_name;type:text;->"`
	Province    *string   `gorm:"column:province;type:text;->"`
	City        *string   `gorm:"column:city;type:text;->"`
	District    string    `gorm:"column:district;type:text;->"`
	Address     string    `gorm:"column:address;type:text;->"`
	Telephone   string    `gorm:"column:telephone;type:varchar(20);->"`
	StarLevel   *int16    `gorm:"column:star_level;type:smallint;->"`
	Description *string   `gorm:"column:description;type:text;->"`
	MainImage   *string   `gorm:"column:main_image;type:text;->"`
	MinPrice    float64   `gorm:"column:min_price;type:numeric;->"`
	RoomCount   int64     `gorm:"column:room_count;type:bigint;->"`
	TotalRooms  int64     `gorm:"column:total_rooms;type:bigint;->"`
	AvgRating   float64   `gorm:"column:avg_rating;type:numeric(3,2);->"`
	ReviewCount int64     `gorm:"column:review_count;type:bigint;->"`
	Status      int16     `gorm:"column:status;type:smallint;->"`
}

func (HotelSummary) TableName() string {
	return "view_hotel_summary_1718"
}
