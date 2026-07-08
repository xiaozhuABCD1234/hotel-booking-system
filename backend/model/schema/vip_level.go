package model

// VipLevel VIP 等级定义，对应表 vip_level_1718
type VipLevel struct {
	Level        int16   `json:"level"        gorm:"column:level;type:smallint;primaryKey;check:level >= 0"`
	LevelName    string  `json:"levelName"    gorm:"column:level_name;type:text;not null"`
	MinPoints    int32   `json:"minPoints"    gorm:"column:min_points;type:integer;not null;check:min_points >= 0"`
	DiscountRate float64 `json:"discountRate" gorm:"column:discount_rate;type:numeric(3,2);not null;check:discount_rate >= 0 AND discount_rate <= 1"`
}

func (VipLevel) TableName() string {
	return "vip_level_1718"
}
