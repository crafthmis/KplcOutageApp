package models

import (
	"time"
)

// Campaign represents the campaign entity
type Campaign struct {
	CmpID       uint `gorm:"primaryKey;column:cmp_id;autoIncrement"`
	OtsID       uint `gorm:"column:ots_id"`
	Msisdn      string
	Message     string
	SentStatus  int       `gorm:"column:sent_status"`
	AckStatus   int       `gorm:"column:ack_status"`
	RecvStatus  int       `gorm:"column:recv_status"`
	DateCreated time.Time `gorm:"column:date_created;autoCreateTime"`
	LastUpdate  time.Time `gorm:"column:last_update;autoUpdateTime"`
	Outage      *Outage   `gorm:"foreignKey:OtsID"`
}

// TableName specifies the table name for the Campaign model
func (Campaign) TableName() string {
	return "tbl_campaign"
}
