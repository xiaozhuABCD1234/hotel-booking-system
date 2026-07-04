package model

import (
	"time"

	"github.com/google/uuid"
)

// Hotel 酒店表，对应表 hotel_1718
type Hotel struct {
	ID          uuid.UUID    `gorm:"column:id;type:uuid;primaryKey;default:uuidv7()"`
	HotelName   string       `gorm:"column:hotel_name;type:text;not null;index:idx_hotel_name"`
	RegionID    int          `gorm:"column:region_id;type:integer;not null;index:idx_hotel_region"`
	Region      Region       `gorm:"foreignKey:RegionID;references:ID"`
	Address     string       `gorm:"column:address;type:text;not null"`
	Telephone   string       `gorm:"column:telephone;type:varchar(20);not null"`
	StarLevel   *int16       `gorm:"column:star_level;type:smallint;check:star_level BETWEEN 1 AND 5;index:idx_hotel_star"`
	Description *string      `gorm:"column:description;type:text"`
	CreateAt    time.Time    `gorm:"column:create_at;type:timestamptz;autoCreateTime"`
	UpdateAt    time.Time    `gorm:"column:update_at;type:timestamptz;autoUpdateTime"`
	Status      int16        `gorm:"column:status;type:smallint;not null;default:1"`
	Images      []HotelImage `gorm:"foreignKey:HotelID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Hotel) TableName() string {
	return "hotel_1718"
}

// HotelImage 酒店图片表，对应表 hotel_image_1718
type HotelImage struct {
	HotelID  uuid.UUID `gorm:"column:hotel_id;type:uuid;primaryKey"`
	ImageURL string    `gorm:"column:image_url;type:text;primaryKey"`
	Hotel    Hotel     `gorm:"foreignKey:HotelID;references:ID;constraint:OnDelete:CASCADE"`
}

func (HotelImage) TableName() string {
	return "hotel_image_1718"
}
