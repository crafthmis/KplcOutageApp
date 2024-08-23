package models

import (
	"time"
)

// Outage represents the outage entity
type Outage struct {
	OtsID       uint `gorm:"primaryKey;column:ots_id;autoIncrement"`
	Message     string
	OutageDate  time.Time    `gorm:"column:outage_date"`
	SentStatus  string       `gorm:"column:sent_status"`
	DateCreated time.Time    `gorm:"column:date_created;autoCreateTime"`
	LastUpdate  time.Time    `gorm:"column:last_update;autoUpdateTime"`
	Areas       []OutageArea `gorm:"foreignKey:OtsID;association_foreignkey:OtsID"`
}

// TableName specifies the table name for the Outage model
func (Outage) TableName() string {
	return "tbl_outage"
}
