package models

import (
	"time"
)

// Outage represents the outage entity
type Outage struct {
	OtsID       uint `gorm:"primaryKey;column:ots_id;autoIncrement"`
	Message     string
	OutageDate  time.Time    `gorm:"column:outage_date"`
	SendStatus  string       `gorm:"column:send_status"`
	Areas       []OutageArea `json:"areas",gorm:"foreignKey:OtsID;association_foreignkey:OtsID"`
	DateCreated time.Time    `gorm:"column:date_created;autoCreateTime"`
	LastUpdate  time.Time    `gorm:"column:last_update;autoUpdateTime"`
}

// TableName specifies the table name for the Outage model
func (Outage) TableName() string {
	return "tbl_outage"
}
