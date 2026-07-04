package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Room 客房表（房型定义），对应表 room_1718
//
// 注意：表级 CHECK (available_quantity <= total_quantity) 在 BeforeSave 中校验。
type Room struct {
	ID                uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:uuidv7()"`
	HotelID           uuid.UUID      `gorm:"column:hotel_id;type:uuid;not null;index:idx_room_hotel"`
	Hotel             Hotel          `gorm:"foreignKey:HotelID;references:ID;constraint:OnDelete:CASCADE"`
	TypeName          string         `gorm:"column:type_name;type:text;not null;index:idx_room_type"`
	TotalQuantity     int32          `gorm:"column:total_quantity;type:integer;not null;check:total_quantity > 0"`
	AvailableQuantity int32          `gorm:"column:available_quantity;type:integer;not null;check:available_quantity >= 0;index:idx_room_available"`
	Price             float64        `gorm:"column:price;type:numeric(10,2);not null;check:price > 0;index:idx_room_price"`
	WeekendPrice      *float64       `gorm:"column:weekend_price;type:numeric(10,2);check:weekend_price > 0"`
	Description       *string        `gorm:"column:description;type:text"`
	CreateAt          time.Time      `gorm:"column:create_at;type:timestamptz;autoCreateTime"`
	UpdateAt          time.Time      `gorm:"column:update_at;type:timestamptz;autoUpdateTime"`
	Status            int16          `gorm:"column:status;type:smallint;not null;default:1"`
	Images            []RoomImage    `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
	Facilities        []RoomFacility `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Room) TableName() string {
	return "room_1718"
}

func (r *Room) BeforeSave(tx *gorm.DB) error {
	if r.AvailableQuantity > r.TotalQuantity {
		return errors.New("available_quantity cannot exceed total_quantity")
	}
	return nil
}

// RoomImage 客房图片表，对应表 room_image_1718
type RoomImage struct {
	RoomID   uuid.UUID `gorm:"column:room_id;type:uuid;primaryKey"`
	ImageURL string    `gorm:"column:image_url;type:text;primaryKey"`
	Room     Room      `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
}

func (RoomImage) TableName() string {
	return "room_image_1718"
}

// RoomFacility 客房设施表，对应表 room_facility_1718
type RoomFacility struct {
	RoomID       uuid.UUID `gorm:"column:room_id;type:uuid;primaryKey"`
	FacilityName string    `gorm:"column:facility_name;type:text;primaryKey"`
	Room         Room      `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
}

func (RoomFacility) TableName() string {
	return "room_facility_1718"
}
