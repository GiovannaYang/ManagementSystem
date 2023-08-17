package model

import "gorm.io/gorm"

type School struct {
	gorm.Model
	Name string `gorm:"type:varchar(50);not null"`
}
