package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Login    string    `gorm:"unique;not null;type:varchar(50)"`
	Password string    `gorm:"not null"`
	Name     string    `gorm:"not null;type:varchar(100)"`
	Email    string    `gorm:"unique;not null;type:varchar(100)"`
	Profiles []Profile `gorm:"many2many:user_profiles"` // Relación muchos a muchos con profiles
}

type Profile struct {
	gorm.Model
	Name        string `gorm:"unique;not null;type:varchar(100)"`
	Description string `gorm:"unique;not null;type:varchar(255)"`
	Users       []User `gorm:"many2many:user_profiles"`
}
