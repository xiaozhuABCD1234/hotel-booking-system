package view

import (
	"time"
)

// PersonInfo 对应视图 view_person_info_1718（只读）
//
// 从身份证号推导性别、年龄、出生日期。
// ⚠️ 被 ViewOrderFull、ViewGuestBookingStats 依赖（SQL 层视图依赖）。
type PersonInfo struct {
	IDCard    string     `json:"idCard"    gorm:"column:id_card;type:char(18);->"`
	Name      string     `json:"name"      gorm:"column:name;type:text;->"`
	Phone     *string    `json:"phone"     gorm:"column:phone;type:varchar(20);->"`
	BirthDate *time.Time `json:"birthDate" gorm:"column:birth_date;type:date;->"`
	Gender    *string    `json:"gender"    gorm:"column:gender;type:text;->"`
	Age       *int       `json:"age"       gorm:"column:age;type:integer;->"`
}

func (PersonInfo) TableName() string {
	return "view_person_info_1718"
}
