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
	ID         uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuidv7()"`
	Username   string    `gorm:"column:username;type:text;not null"`
	Password   string    `gorm:"column:password;type:varchar(255);not null"`
	Phone      *string   `gorm:"column:phone;type:varchar(20);index:idx_user_phone"`
	Email      *string   `gorm:"column:email;type:text"`
	RealName   *string   `gorm:"column:real_name;type:text"`
	IDCard     *string   `gorm:"column:id_card;type:varchar(18);check:(id_card ~ '^\\d{17}[\\dXx]$' OR id_card IS NULL)"`
	Role       UserRole  `gorm:"column:role;type:user_role;not null;default:customer;index:idx_user_role"`
	Points     int32     `gorm:"column:points;type:integer;not null;default:0;check:points >= 0;index:idx_user_points,sort:desc"`
	VipLevelID int16     `gorm:"column:vip_level;type:smallint;not null;default:0"`
	VipLevel   VipLevel  `gorm:"foreignKey:VipLevelID;references:Level"`
	CreateAt   time.Time `gorm:"column:create_at;type:timestamptz;autoCreateTime"`
	UpdateAt   time.Time `gorm:"column:update_at;type:timestamptz;autoUpdateTime"`
	Status     int16     `gorm:"column:status;type:smallint;not null;default:1"`
}

func (User) TableName() string {
	return "user_1718"
}
