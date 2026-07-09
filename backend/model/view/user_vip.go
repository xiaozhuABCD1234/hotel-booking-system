package view

import (
	"time"

	"github.com/google/uuid"
)

// UserVip 对应视图 view_user_vip_1718（只读）
//
// 用户 VIP 信息，用于个人中心、下单折扣计算。
// Role 以 string 接收 user_role 枚举值，避免与 schema 包产生跨包依赖。
type UserVip struct {
	UserID            uuid.UUID `json:"userId"            gorm:"column:user_id;type:uuid;->"`
	Username          string    `json:"username"          gorm:"column:username;type:text;->"`
	Phone             *string   `json:"phone"             gorm:"column:phone;type:varchar(20);->"`
	Email             *string   `json:"email"             gorm:"column:email;type:text;->"`
	RealName          *string   `json:"realName"          gorm:"column:real_name;type:text;->"`
	IDCard            *string   `json:"idCard"            gorm:"column:id_card;type:varchar(18);->"`
	Role              string    `json:"role"              gorm:"column:role;type:user_role;->"`
	Points            int32     `json:"points"            gorm:"column:points;type:integer;->"`
	VipLevel          int16     `json:"vipLevel"          gorm:"column:vip_level;type:smallint;->"`
	VipLevelName      string    `json:"vipLevelName"      gorm:"column:vip_level_name;type:text;->"`
	DiscountRate      float64   `json:"discountRate"      gorm:"column:discount_rate;type:numeric(3,2);->"`
	PointsToNextLevel *int32    `json:"pointsToNextLevel" gorm:"column:points_to_next_level;type:integer;->"`
	CreateAt          time.Time `json:"createAt"          gorm:"column:create_at;type:timestamptz;->"`
}

func (UserVip) TableName() string {
	return "view_user_vip_1718"
}
