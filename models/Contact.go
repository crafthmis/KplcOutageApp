package models

import (
	"time"

	"gorm.io/gorm"
)

// Contact represents the contact entity
type Contact struct {
	CntID        uint `gorm:"primaryKey;column:cnt_id;autoIncrement"`
	AreaID       uint `gorm:"column:area_id"`
	SubID        int  `gorm:"column:sub_id"`
	Msisdn       string
	FirstName    string        `gorm:"column:first_name"`
	LastName     string        `gorm:"column:last_name"`
	Username     string        `json:"-",gorm:"column:username"`
	Password     string        `gorm:"column:password"`
	Email        string        `gorm:"column:email"`
	AuthToken    string        `json:"-",gorm:"column:auth_token"`
	DateCreated  time.Time     `json:"-",gorm:"column:date_created;autoCreateTime"`
	LastUpdate   time.Time     `json:"-",gorm:"column:last_update;autoUpdateTime"`
	Area         *Area         `gorm:"foreignKey:AreaID"`
	Subscription *Subscription `gorm:"belongsTo:Subscription;foreignKey:SubID"`
}

// TableName specifies the table name for the Contact model
func (Contact) TableName() string {
	return "tbl_contact"
}

// Example of BeforeCreate hook

func (user *Contact) BeforeCreate(tx *gorm.DB) (err error) {
	s := user.Msisdn

	user.Msisdn = s[len(s)-9:]

	return nil

}
