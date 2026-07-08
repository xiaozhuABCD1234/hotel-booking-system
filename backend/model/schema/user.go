package model

import (
	"time"

	"github.com/google/uuid"
)

// UserRole 用户角色枚举，对应 PostgreSQL 自定义类型 user_role。
type UserRole string

const (
	RoleCustomer     UserRole = "customer"
	RoleVIP          UserRole = "vip"
	RoleHotelManager UserRole = "hotel_manager"
	RoleAdmin        UserRole = "admin"
)

// User 用户表，对应表 user_1718
type User struct {
	ID         uuid.UUID `json:"id"         gorm:"column:id;type:uuid;primaryKey;default:uuidv7()"`
	Username   string    `json:"username"   gorm:"column:username;type:text;not null"`
	Password   string    `json:"-"          gorm:"column:password;type:varchar(255);not null"`
	Phone      *string   `json:"phone"      gorm:"column:phone;type:varchar(20);index:idx_user_phone"`
	Email      *string   `json:"email"      gorm:"column:email;type:text"`
	RealName   *string   `json:"realName"   gorm:"column:real_name;type:text"`
	IDCard     *string   `json:"idCard"     gorm:"column:id_card;type:varchar(18);check:(id_card ~ '^\\d{17}[\\dXx]$' OR id_card IS NULL)"`
	Role       UserRole  `json:"role"       gorm:"column:role;type:user_role;not null;default:customer;index:idx_user_role"`
	Points     int32     `json:"points"     gorm:"column:points;type:integer;not null;default:0;check:points >= 0;index:idx_user_points,sort:desc"`
	VipLevelID int16     `json:"vipLevelId" gorm:"column:vip_level;type:smallint;not null;default:0"`
	VipLevel   VipLevel  `json:"vipLevel"   gorm:"foreignKey:VipLevelID;references:Level"`
	CreateAt   time.Time `json:"createAt"   gorm:"column:create_at;type:timestamptz;autoCreateTime"`
	UpdateAt   time.Time `json:"updateAt"   gorm:"column:update_at;type:timestamptz;autoUpdateTime"`
	Status     int16     `json:"status"     gorm:"column:status;type:smallint;not null;default:1"`
}

func (User) TableName() string {
	return "user_1718"
}
