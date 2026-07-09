package model

// Person 入住人员身份信息表，对应表 person_1718
type Person struct {
	IDCard     string  `json:"idCard"     gorm:"column:id_card;type:char(18);primaryKey;check:(id_card ~ '^\\d{17}[\\dXx]$')"`
	Name       string  `json:"name"       gorm:"column:name;type:text;not null"`
	Phone      *string `json:"phone"      gorm:"column:phone;type:varchar(20)"`
	Occupation *string `json:"occupation" gorm:"column:occupation;type:text"`
	Education  *string `json:"education"  gorm:"column:education;type:education_level"`
	Income     *string `json:"income"     gorm:"column:income;type:numrange"`
}

func (Person) TableName() string {
	return "person_1718"
}
