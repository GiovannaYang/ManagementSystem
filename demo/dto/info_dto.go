package dto

import (
	"demo/common"
	"demo/model"
)

type InfoDto struct {
	ID        uint   `json:"id"`
	Sid       string `json:"sid"`
	School    string `json:"school"`
	Major     string `json:"major"`
	Class     string `json:"class""`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}

func ToInfoDto(info model.Info) InfoDto {
	db := common.GetDB()
	var school string
	db.Raw("SELECT schools.Name FROM schools,classes WHERE schools.ID=classes.school_id AND classes.ID= ?", info.ClassID).Scan(&school)
	var major string
	db.Raw("SELECT majors.Name FROM majors,classes WHERE majors.ID=classes.major_id AND classes.ID = ?", info.ClassID).Scan(&major)
	var class string
	db.Raw("SELECT classes.Name FROM classes WHERE classes.ID=?", info.ClassID).Scan(&class)
	return InfoDto{
		ID:        info.ID,
		Sid:       info.Sid,
		School:    school,
		Major:     major,
		Class:     class,
		Name:      info.Name,
		Gender:    info.Gender,
		Email:     info.Email,
		Telephone: info.Telephone,
	}
}
