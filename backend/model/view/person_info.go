package view

import (
	"time"
)

// PersonInfo 对应视图 view_person_info_1718（只读）
//
// 从身份证号推导性别、年龄、出生日期。
// ⚠️ 被 ViewOrderFull、ViewGuestBookingStats 依赖（SQL 层视图依赖）。
type PersonInfo struct {
	IDCard    string     `gorm:"column:id_card;type:char(18);->"`
	Name      string     `gorm:"column:name;type:text;->"`
	Phone     *string    `gorm:"column:phone;type:varchar(20);->"`
	BirthDate *time.Time `gorm:"column:birth_date;type:date;->"`
	Gender    *string    `gorm:"column:gender;type:text;->"`
	Age       *int       `gorm:"column:age;type:integer;->"`
}

func (PersonInfo) TableName() string {
	return "view_person_info_1718"
}
