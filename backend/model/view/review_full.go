package view

import (
	"time"

	"github.com/google/uuid"
)

// ReviewFull 对应视图 view_review_full_1718（只读）
//
// 评价详情，用于酒店评价列表页、用户评价记录。
type ReviewFull struct {
	ReviewID     uuid.UUID `json:"reviewId"     gorm:"column:review_id;type:uuid;->"`
	UserID       uuid.UUID `json:"userId"       gorm:"column:user_id;type:uuid;->"`
	Username     string    `json:"username"     gorm:"column:username;type:text;->"`
	HotelID      uuid.UUID `json:"hotelId"      gorm:"column:hotel_id;type:uuid;->"`
	HotelName    string    `json:"hotelName"    gorm:"column:hotel_name;type:text;->"`
	OrderID      uuid.UUID `json:"orderId"      gorm:"column:order_id;type:uuid;->"`
	RoomType     string    `json:"roomType"     gorm:"column:room_type;type:text;->"`
	CheckInDate  time.Time `json:"checkInDate"  gorm:"column:check_in_date;type:date;->"`
	CheckOutDate time.Time `json:"checkOutDate" gorm:"column:check_out_date;type:date;->"`
	Rating       int16     `json:"rating"       gorm:"column:rating;type:smallint;->"`
	Content      *string   `json:"content"      gorm:"column:content;type:text;->"`
	CreateAt     time.Time `json:"createAt"     gorm:"column:create_at;type:timestamptz;->"`
}

func (ReviewFull) TableName() string {
	return "view_review_full_1718"
}
