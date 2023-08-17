package controller

import (
	"demo/common"
	"demo/model"
	"demo/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Info2 struct {
	ClassName  string
	SchoolName string
	MajorName  string
}

func AddClass(ctx *gin.Context) {
	db := common.GetDB()
	//获取参数
	var requestClass model.Class
	ctx.Bind(&requestClass)
	name := requestClass.Name

	if isClassExist(db, name) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}
	newClass := model.Class{
		Name: name,
	}
	db.Create(&newClass)

	response.Success(ctx, nil, "创建成功")
}

func ShowClass(ctx *gin.Context) {
	db := common.GetDB()
	var classInfo []Info2
	db.Raw("SELECT schools.name as SchoolName,majors.name as MajorName,classes.name as ClassName FROM classes LEFT JOIN majors ON classes.major_id=majors.ID LEFT JOIN schools ON classes.school_id=schools.ID").Scan(&classInfo)
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": classInfo})
}

func DeleteClass(ctx *gin.Context) {
	name := ctx.Params.ByName("name")
	var cla model.Class

	db := common.GetDB()
	db.AutoMigrate(model.Class{})
	db.Where("name = ?", name).Unscoped().Delete(&cla)
	response.Success(ctx, nil, "删除成功")
}

func UpdateClass(ctx *gin.Context) {
	// 绑定body中的参数
	var requestClass Info2
	ctx.Bind(&requestClass)
	majorName := requestClass.MajorName
	schoolName := requestClass.SchoolName

	name := ctx.Params.ByName("name")

	db := common.GetDB()
	db.AutoMigrate(model.Class{})
	var cla model.Class
	db.Where("name = ?", name).First(&cla)

	updateClass := model.Class{
		Name:     "",
		MajorID:  0,
		SchoolID: 0,
	}

	if schoolName != "" {
		var school model.School
		db.Where("name = ?", schoolName).First(&school)
		if school.ID == 0 {
			response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "该学院不存在")
			return
		}
		updateClass.SchoolID = school.ID
	}

	if majorName != "" {
		var major model.Major
		db.Where("name = ?", majorName).First(&major)
		if major.ID == 0 {
			response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "该专业不存在")
			return
		}
		updateClass.MajorID = major.ID
	}
	db.Model(&cla).Updates(updateClass)
	db.Raw("SELECT * FROM classes WHERE name = ?", name).Scan(&updateClass)
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"class": updateClass}})
}

func isClassExist(db *gorm.DB, name string) bool {
	var cla model.Class
	db.Where("name = ?", name).First(&cla)
	if cla.ID != 0 {
		return true
	}
	return false
}
