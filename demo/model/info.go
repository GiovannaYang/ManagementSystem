package model

type Info struct {
	ID        uint   `gorm:"primary"`
	Sid       string `gorm:"type:varchar(20);not null;unique"`
	ClassID   uint   `gorm:"default:0"`
	Name      string `gorm:"type:varchar(20);default:''"`
	Gender    string `gorm:"type:varchar(2);default:''"`
	Email     string `gorm:"size:50;default:''"`
	Telephone string `gorm:"type:varchar(11);default:''"`
}
