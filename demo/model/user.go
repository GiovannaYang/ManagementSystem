package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Sid      string `gorm:"type:varchar(20);not null;unique"`
	Password string `gorm:"size:255;not null"`
}
