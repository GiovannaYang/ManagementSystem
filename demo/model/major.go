package model

import "gorm.io/gorm"

type Major struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20);not null"`
	SchoolID uint   `gorm:"not null"`
}
