package model

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Email    string `gorm:"size:50;not null;unique"`
	Name     string `gorm:"type:varchar(20);not null"`
	Password string `gorm:"size:255;not null"`
}
