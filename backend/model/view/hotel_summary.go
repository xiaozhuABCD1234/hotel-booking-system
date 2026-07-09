package view

import (
	"github.com/google/uuid"
)

// HotelSummary 对应视图 view_hotel_summary_1718（只读）
//
// 酒店摘要列表，用于搜索页、首页推荐。
// min_price / avg_rating 经 COALESCE 兜底为 0，不会为 NULL。
type HotelSummary struct {
	HotelID     uuid.UUID `json:"hotelId"     gorm:"column:hotel_id;type:uuid;->"`
	HotelName   string    `json:"hotelName"   gorm:"column:hotel_name;type:text;->"`
	Province    *string   `json:"province"    gorm:"column:province;type:text;->"`
	City        *string   `json:"city"        gorm:"column:city;type:text;->"`
	District    string    `json:"district"    gorm:"column:district;type:text;->"`
	Address     string    `json:"address"     gorm:"column:address;type:text;->"`
	Telephone   string    `json:"telephone"   gorm:"column:telephone;type:varchar(20);->"`
	StarLevel   *int16    `json:"starLevel"   gorm:"column:star_level;type:smallint;->"`
	Description *string   `json:"description" gorm:"column:description;type:text;->"`
	MainImage   *string   `json:"mainImage"   gorm:"column:main_image;type:text;->"`
	MinPrice    float64   `json:"minPrice"    gorm:"column:min_price;type:numeric;->"`
	RoomCount   int64     `json:"roomCount"   gorm:"column:room_count;type:bigint;->"`
	TotalRooms  int64     `json:"totalRooms"  gorm:"column:total_rooms;type:bigint;->"`
	AvgRating   float64   `json:"avgRating"   gorm:"column:avg_rating;type:numeric(3,2);->"`
	ReviewCount int64     `json:"reviewCount" gorm:"column:review_count;type:bigint;->"`
	Status      int16     `json:"status"      gorm:"column:status;type:smallint;->"`
}

func (HotelSummary) TableName() string {
	return "view_hotel_summary_1718"
}
