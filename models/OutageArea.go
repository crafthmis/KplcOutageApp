package models

// OutageArea represents the outage area entity
type OutageArea struct {
	OctID  uint    `gorm:"primaryKey;column:oct_id;autoIncrement"`
	AreaID uint    `gorm:"column:area_id"`
	OtsID  uint    `gorm:"column:ots_id"`
	Outage *Outage `gorm:"foreignKey:OtsID"`
}

// TableName specifies the table name for the OutageArea model
func (OutageArea) TableName() string {
	return "tbl_outage_area"
}
