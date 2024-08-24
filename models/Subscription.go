package models

import "time"

type Subscription struct {
	SubID       uint      `gorm:"column:sub_id;primaryKey;autoIncrement"`
	Name        string    `gorm:"column:name;not null"`
	Description string    `gorm:"column:description;not null"`
	Amount      int64     `gorm:"column:amount;not null"`
	DateCreated time.Time `json:"-" , gorm:"column:date_created;autoCreateTime"`
	LastUpdate  time.Time `json:"-" , gorm:"column:last_update;autoUpdateTime"`
	Contact     *Contact  `gorm:"foreignKey:SubID;association_foreignkey:SubID"`
}

func (Subscription) TableName() string {
	return "tbl_subscription"
}
