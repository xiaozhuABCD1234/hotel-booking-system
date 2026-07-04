package view

import (
	"time"

	"github.com/google/uuid"
)

// ReviewFull 对应视图 view_review_full_1718（只读）
//
// 评价详情，用于酒店评价列表页、用户评价记录。
type ReviewFull struct {
	ReviewID     uuid.UUID `gorm:"column:review_id;type:uuid;->"`
	UserID       uuid.UUID `gorm:"column:user_id;type:uuid;->"`
	Username     string    `gorm:"column:username;type:text;->"`
	HotelID      uuid.UUID `gorm:"column:hotel_id;type:uuid;->"`
	HotelName    string    `gorm:"column:hotel_name;type:text;->"`
	OrderID      uuid.UUID `gorm:"column:order_id;type:uuid;->"`
	RoomType     string    `gorm:"column:room_type;type:text;->"`
	CheckInDate  time.Time `gorm:"column:check_in_date;type:date;->"`
	CheckOutDate time.Time `gorm:"column:check_out_date;type:date;->"`
	Rating       int16     `gorm:"column:rating;type:smallint;->"`
	Content      *string   `gorm:"column:content;type:text;->"`
	CreateAt     time.Time `gorm:"column:create_at;type:timestamptz;->"`
}

func (ReviewFull) TableName() string {
	return "view_review_full_1718"
}
