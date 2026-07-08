package model

import (
	"time"

	"github.com/google/uuid"
)

// Hotel 酒店表，对应表 hotel_1718
type Hotel struct {
	ID          uuid.UUID    `json:"id"          gorm:"column:id;type:uuid;primaryKey;default:uuidv7()"`
	HotelName   string       `json:"hotelName"   gorm:"column:hotel_name;type:text;not null;index:idx_hotel_name"`
	RegionID    int          `json:"regionId"    gorm:"column:region_id;type:integer;not null;index:idx_hotel_region"`
	Region      Region       `json:"region"      gorm:"foreignKey:RegionID;references:ID"`
	Address     string       `json:"address"     gorm:"column:address;type:text;not null"`
	Telephone   string       `json:"telephone"   gorm:"column:telephone;type:varchar(20);not null"`
	StarLevel   *int16       `json:"starLevel"   gorm:"column:star_level;type:smallint;check:star_level BETWEEN 1 AND 5;index:idx_hotel_star"`
	Description *string      `json:"description" gorm:"column:description;type:text"`
	CreateAt    time.Time    `json:"createAt"    gorm:"column:create_at;type:timestamptz;autoCreateTime"`
	UpdateAt    time.Time    `json:"updateAt"    gorm:"column:update_at;type:timestamptz;autoUpdateTime"`
	Status      int16        `json:"status"      gorm:"column:status;type:smallint;not null;default:1"`
	Images      []HotelImage `json:"images"      gorm:"foreignKey:HotelID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Hotel) TableName() string {
	return "hotel_1718"
}

// HotelImage 酒店图片表，对应表 hotel_image_1718
type HotelImage struct {
	HotelID  uuid.UUID `json:"hotelId"  gorm:"column:hotel_id;type:uuid;primaryKey"`
	ImageURL string    `json:"imageUrl" gorm:"column:image_url;type:text;primaryKey"`
	Hotel    Hotel     `json:"-"        gorm:"foreignKey:HotelID;references:ID;constraint:OnDelete:CASCADE"`
}

func (HotelImage) TableName() string {
	return "hotel_image_1718"
}
