package models

// OutageArea represents the outage area entity
type OutageArea struct {
	OctID  uint    `gorm:"primaryKey;column:oct_id;autoIncrement"`
	AreaID uint    `gorm:"column:area_id",gorm:"index:tbl_outage_area_ots_id_area_id_key,unique"`
	OtsID  uint    `gorm:"column:ots_id",gorm:"index:tbl_outage_area_ots_id_area_id_key,unique"`
	Outage *Outage `gorm:"foreignKey:OtsID"`
}

// TableName specifies the table name for the OutageArea model
func (OutageArea) TableName() string {
	return "tbl_outage_area"
}
