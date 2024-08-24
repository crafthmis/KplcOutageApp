package models

import (
	"time"
)

// Contact represents the contact entity
type Contact struct {
	CntID       uint `gorm:"primaryKey;column:cnt_id;autoIncrement"`
	AreaID      uint `gorm:"column:area_id"`
	SubID       int  `gorm:"column:sub_id"`
	Msisdn      string
	FirstName   string        `gorm:"column:first_name"`
	LastName    string        `gorm:"column:last_name"`
	UserName    string        `json:"-",gorm:"column:username"`
	Password    string        `json:"-",gorm:"column:password"`
	Email       string        `json:"-",gorm:"column:email"`
	Token       string        `json:"-",gorm:"column:auth_token"`
	DateCreated time.Time     `json:"-",gorm:"column:date_created;autoCreateTime"`
	LastUpdate  time.Time     `json:"-",gorm:"column:last_update;autoUpdateTime"`
	Area        *Area         `gorm:"foreignKey:AreaID"`
	Role        *Subscription `gorm:"foreignKey:SubID;association_foreignkey:SubID"`

}

// TableName specifies the table name for the Contact model
func (Contact) TableName() string {
	return "tbl_contact"
}
