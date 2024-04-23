package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Login    string    `gorm:"unique;not null;type:varchar(50)"`
	Password string    `gorm:"not null"`
	Name     string    `gorm:"not null;type:varchar(100)"`
	Email    string    `gorm:"unique;not null;type:varchar(100)"`
	Profiles []Profile `gorm:"many2many:sec_users_profiles"`
}

type Profile struct {
	gorm.Model
	Description string `gorm:"unique;not null;type:varchar(255)"`
	PrivAccess  bool   `gorm:"not null;default:false"`
	PrivInsert  bool   `gorm:"not null;default:false"`
	PrivUpdate  bool   `gorm:"not null;default:false"`
	PrivDelete  bool   `gorm:"not null;default:false"`
	PrivExport  bool   `gorm:"not null;default:false"`
	PrivPrint   bool   `gorm:"not null;default:false"`
	Users       []User `gorm:"many2many:sec_users_profiles"`
}

func (User) TableName() string {
	return "sec_users"
}

func (Profile) TableName() string {
	return "sec_profiles"
}
