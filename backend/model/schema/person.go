package model

// Person 入住人员身份信息表，对应表 person_1718
type Person struct {
	IDCard string  `gorm:"column:id_card;type:varchar(18);primaryKey;check:(id_card ~ '^\\d{17}[\\dXx]$')"`
	Name   string  `gorm:"column:name;type:text;not null"`
	Phone  *string `gorm:"column:phone;type:varchar(20)"`
}

func (Person) TableName() string {
	return "person_1718"
}
