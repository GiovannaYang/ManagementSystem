package controller

import (
	"demo/common"
	"demo/dto"
	"demo/model"
	"demo/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type Info3 struct {
	ID        uint
	Sid       string
	ClassName string
	Name      string
	Gender    string
	Email     string
	Telephone string
}

func Register(ctx *gin.Context) {
	db := common.GetDB()
	//获取参数
	var requestUser = model.User{}
	ctx.Bind(&requestUser)
	sid := requestUser.Sid
	password := requestUser.Password
	//数据验证
	if len(password) < 6 && password != "" {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	if password == "" {
		password = "123456"
	} // 如果进行了设置密码长度不小于6，如果不进行设置为默认123456

	if sid == "" {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "请输入学号")
		return
	}

	//判断学号是否存在
	if isSidExist(db, sid) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}

	//创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	newUser := model.User{
		Sid:      sid,
		Password: string(hasedPassword),
	}
	db.Create(&newUser)
	newInfo := model.Info{
		ID:  newUser.ID,
		Sid: sid,
	}
	db.Create(&newInfo)
	//返回结果
	response.Success(ctx, nil, "注册成功")
}

func Login(ctx *gin.Context) {
	db := common.GetDB()
	//获取参数
	var requestUser = model.User{}
	ctx.Bind(&requestUser)
	sid := requestUser.Sid
	password := requestUser.Password

	// 数据验证
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	// 判断学号是否存在
	var user model.User
	db.Where("sid = ?", sid).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "密码错误")
		return
	}
	// 发放token
	token, err := common.ReleaseToken(user.ID)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	// 返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

func Info(ctx *gin.Context) {
	info, _ := ctx.Get("data")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToInfoDto(info.(model.Info))}})
}

func PassUpdate(ctx *gin.Context) {
	// 绑定body中的参数
	var requestUser model.User
	ctx.Bind(&requestUser)
	password := requestUser.Password

	userId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	var user model.User
	user.ID = uint(userId)

	db := common.GetDB()
	db.AutoMigrate(model.User{})

	updateUser := model.User{
		Password: "",
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
		updateUser.Password = string(hasedPassword)
	}
	db.Model(&user).Updates(updateUser)
	response.Success(ctx, gin.H{"user": user}, "修改成功")
}

func InfoUpdate(ctx *gin.Context) {
	// 绑定body中的参数
	var requestInfo Info3
	ctx.Bind(&requestInfo)
	name := requestInfo.Name
	className := requestInfo.ClassName
	gender := requestInfo.Gender
	email := requestInfo.Email
	telephone := requestInfo.Telephone

	userId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	var info model.Info
	info.ID = uint(userId)

	db := common.GetDB()
	db.AutoMigrate(model.Info{})

	updateInfo := model.Info{
		Name:      "",
		ClassID:   0,
		Gender:    "",
		Email:     "",
		Telephone: "",
	}
	if className != "" {
		var class model.Class
		db.Where("name = ?", className).First(&class)
		if class.ID == 0 {
			response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "该班级不存在")
			return
		}
		updateInfo.ClassID = class.ID
	}
	if name != "" {
		updateInfo.Name = name
	}
	if gender != "" {
		if gender != "男" && gender != "女" {
			response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "性别必须为男或女")
			return
		}
		updateInfo.Gender = gender
	}
	if email != "" {
		updateInfo.Email = email
	}
	if telephone != "" {
		if len(telephone) != 11 {
			response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "电话号码必须为11位")
			return
		}
		updateInfo.Telephone = telephone
	}
	db.Model(&info).Updates(updateInfo)
	db.Raw("SELECT * FROM infos WHERE id = ?", userId).Scan(&updateInfo)
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToInfoDto(updateInfo)}})
}

func Delete(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	var user model.User
	var userinfo model.Info

	db := common.GetDB()
	db.AutoMigrate(model.User{})
	db.AutoMigrate(model.Info{})
	db.Unscoped().Delete(&user, userId)
	db.Unscoped().Delete(&userinfo, userId)
	response.Success(ctx, nil, "")
}

func isSidExist(db *gorm.DB, sid string) bool {
	var user model.User
	db.Where("sid = ?", sid).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
