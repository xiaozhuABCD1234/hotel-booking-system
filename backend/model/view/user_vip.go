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
	UserID            uuid.UUID `gorm:"column:user_id;type:uuid;->"`
	Username          string    `gorm:"column:username;type:text;->"`
	Phone             *string   `gorm:"column:phone;type:varchar(20);->"`
	Email             *string   `gorm:"column:email;type:text;->"`
	RealName          *string   `gorm:"column:real_name;type:text;->"`
	IDCard            *string   `gorm:"column:id_card;type:varchar(18);->"`
	Role              string    `gorm:"column:role;type:user_role;->"`
	Points            int32     `gorm:"column:points;type:integer;->"`
	VipLevel          int16     `gorm:"column:vip_level;type:smallint;->"`
	VipLevelName      string    `gorm:"column:vip_level_name;type:text;->"`
	DiscountRate      float64   `gorm:"column:discount_rate;type:numeric(3,2);->"`
	PointsToNextLevel *int32    `gorm:"column:points_to_next_level;type:integer;->"`
	CreateAt          time.Time `gorm:"column:create_at;type:timestamptz;->"`
}

func (UserVip) TableName() string {
	return "view_user_vip_1718"
}
