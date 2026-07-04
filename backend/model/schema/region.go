package model

// Region 地区表（省/市/区 三级层次），对应表 region_1718
type Region struct {
	ID         int     `gorm:"column:id;type:serial;primaryKey;autoIncrement"`
	RegionName string  `gorm:"column:region_name;type:text;not null"`
	ParentsID  *int    `gorm:"column:parents_id;type:integer;index:idx_region_parents"`
	Parent     *Region `gorm:"foreignKey:ParentsID;references:ID"`
}

func (Region) TableName() string {
	return "region_1718"
}
