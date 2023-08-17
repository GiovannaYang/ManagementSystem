package model

import "gorm.io/gorm"

type Msg struct {
	gorm.Model
	Sid string `gorm:"type:varchar(20);not null"`
}
