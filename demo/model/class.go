package model

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20);not null"`
	MajorID  uint   `gorm:""`
	SchoolID uint   `gorm:""`
}
