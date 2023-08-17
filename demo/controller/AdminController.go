package controller

import (
	"demo/common"
	"demo/dto"
	"demo/model"
	"demo/response"
	"demo/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type Info1 struct {
	ID         uint
	Sid        string
	SchoolName string
	MajorName  string
	ClassName  string
	Name       string
	Gender     string
	Email      string
	Telephone  string
}

func ARegister(ctx *gin.Context) {
	db := common.GetDB()
	//获取参数
	var requestAdmin = model.Admin{}
	err := ctx.Bind(&requestAdmin)
	if err != nil {
		return
	}
	email := requestAdmin.Email
	name := requestAdmin.Name
	password := requestAdmin.Password
	//数据验证
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	// 如果名称没有传，给一个10位大随机字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	//判断邮箱是否存在
	if isEmailExist(db, email) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}
	//创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	newAdmin := model.Admin{
		Email:    email,
		Name:     name,
		Password: string(hasedPassword),
	}
	db.Create(&newAdmin)

	// 发放token
	token, err := common.ReleaseToken(newAdmin.ID)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	// 返回结果
	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

func ALogin(ctx *gin.Context) {
	db := common.GetDB()
	//获取参数
	var requestAdmin = model.Admin{}
	ctx.Bind(&requestAdmin)
	email := requestAdmin.Email
	password := requestAdmin.Password

	// 数据验证
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	// 判断邮箱是否存在
	var admin model.Admin
	db.Where("email = ?", email).First(&admin)
	if admin.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "密码错误")
		return
	}
	// 发放token
	token, err := common.ReleaseToken(admin.ID)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	// 返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

func AdminInfo(ctx *gin.Context) {
	admin, _ := ctx.Get("admin")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"admin": dto.ToAdminDto(admin.(model.Admin))}})
}

func AUpdate(ctx *gin.Context) {
	// 绑定body中的参数
	var requestAdmin model.Admin
	ctx.Bind(&requestAdmin)
	name := requestAdmin.Name
	password := requestAdmin.Password

	adminId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	var admin model.Admin
	admin.ID = uint(adminId)

	db := common.GetDB()
	db.AutoMigrate(model.Admin{})

	updateAdmin := model.Admin{
		Email:    "",
		Name:     "",
		Password: "",
	}
	if name != "" {
		updateAdmin.Name = name
	}
	if password != "" {
		//数据验证
		if len(password) < 6 {
			response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
			return
		}
		hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
			return
		}
		updateAdmin.Password = string(hasedPassword)
	}
	db.Model(&admin).Updates(updateAdmin)
	response.Success(ctx, gin.H{"admin": admin}, "修改成功")
}

func ADelete(ctx *gin.Context) {
	adminId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	var admin model.Admin
	db := common.GetDB()
	db.AutoMigrate(model.Admin{})
	db.Unscoped().Delete(&admin, adminId)
	response.Success(ctx, nil, "")
}

func AShow(ctx *gin.Context) {
	db := common.GetDB()
	var userinfo []Info1
	db.Raw("SELECT infos.ID,infos.sid,schools.name as SchoolName,majors.name as MajorName,classes.name as ClassName,infos.name,infos.gender,infos.email,infos.telephone FROM infos LEFT JOIN classes ON infos.class_id=classes.ID LEFT JOIN majors ON classes.major_id=majors.ID LEFT JOIN schools ON majors.school_id=schools.ID").Scan(&userinfo)
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": userinfo})
}

func ASearch(ctx *gin.Context) {
	word := ctx.Query("word")
	if word == "" {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "查找字段不能为空")
		return
	}

	db := common.GetDB()
	var userinfo []Info1
	db.Raw("SELECT infos.ID,infos.sid,schools.name as SchoolName,majors.name as MajorName,classes.name as ClassName,infos.name,infos.gender,infos.email,infos.telephone FROM infos LEFT JOIN classes ON infos.class_id=classes.ID LEFT JOIN majors ON classes.major_id=majors.ID LEFT JOIN schools ON majors.school_id=schools.ID WHERE sid LIKE ?", "%"+word+"%").Scan(&userinfo)
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": userinfo})
}

func isEmailExist(db *gorm.DB, email string) bool {
	var admin model.Admin
	db.Where("email = ?", email).First(&admin)
	if admin.ID != 0 {
		return true
	}
	return false
}
