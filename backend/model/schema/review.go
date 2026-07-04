package model

import (
	"time"

	"github.com/google/uuid"
)

// Review 评价表，对应表 review_1718
//
// 唯一约束 (user_id, order_id) 保证每笔订单每个用户只能评价一次。
type Review struct {
	ID       uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuidv7()"`
	UserID   uuid.UUID `gorm:"column:user_id;type:uuid;not null;uniqueIndex:idx_review_user_order;index:idx_review_user"`
	User     User      `gorm:"foreignKey:UserID;references:ID"`
	OrderID  uuid.UUID `gorm:"column:order_id;type:uuid;not null;uniqueIndex:idx_review_user_order"`
	Order    Order     `gorm:"foreignKey:OrderID;references:ID"`
	HotelID  uuid.UUID `gorm:"column:hotel_id;type:uuid;not null;index:idx_review_hotel"`
	Hotel    Hotel     `gorm:"foreignKey:HotelID;references:ID"`
	Rating   int16     `gorm:"column:rating;type:smallint;not null;check:rating BETWEEN 1 AND 5;index:idx_review_rating"`
	Content  *string   `gorm:"column:content;type:text"`
	CreateAt time.Time `gorm:"column:create_at;type:timestamptz;autoCreateTime"`
	UpdateAt time.Time `gorm:"column:update_at;type:timestamptz;autoUpdateTime"`
}

func (Review) TableName() string {
	return "review_1718"
}
