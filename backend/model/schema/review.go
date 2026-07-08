package model

import (
	"time"

	"github.com/google/uuid"
)

// Review 评价表，对应表 review_1718
//
// 唯一约束 (user_id, order_id) 保证每笔订单每个用户只能评价一次。
type Review struct {
	ID       uuid.UUID `json:"id"       gorm:"column:id;type:uuid;primaryKey;default:uuidv7()"`
	UserID   uuid.UUID `json:"userId"   gorm:"column:user_id;type:uuid;not null;uniqueIndex:idx_review_user_order;index:idx_review_user"`
	User     User      `json:"user"     gorm:"foreignKey:UserID;references:ID"`
	OrderID  uuid.UUID `json:"orderId"  gorm:"column:order_id;type:uuid;not null;uniqueIndex:idx_review_user_order"`
	Order    Order     `json:"order"    gorm:"foreignKey:OrderID;references:ID"`
	HotelID  uuid.UUID `json:"hotelId"  gorm:"column:hotel_id;type:uuid;not null;index:idx_review_hotel"`
	Hotel    Hotel     `json:"hotel"    gorm:"foreignKey:HotelID;references:ID"`
	Rating   int16     `json:"rating"   gorm:"column:rating;type:smallint;not null;check:rating BETWEEN 1 AND 5;index:idx_review_rating"`
	Content  *string   `json:"content"  gorm:"column:content;type:text"`
	CreateAt time.Time `json:"createAt" gorm:"column:create_at;type:timestamptz;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"column:update_at;type:timestamptz;autoUpdateTime"`
}

func (Review) TableName() string {
	return "review_1718"
}
