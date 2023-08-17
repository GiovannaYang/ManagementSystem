package controller

import (
	"demo/common"
	"demo/model"
	"demo/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

func ForgetPass(ctx *gin.Context) {
	db := common.GetDB()
	var requestUser = model.Msg{}
	ctx.Bind(&requestUser)
	sid := requestUser.Sid
	if isExist(db, sid) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "申请尚未处理")
		return
	}
	if !isExist2(db, sid) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "学号不存在")
		return
	}
	newMessage := model.Msg{
		Sid: sid,
	}
	db.Create(&newMessage)
	// 返回结果
	response.Success(ctx, nil, "申请成功")
}

func ShowMsg(ctx *gin.Context) {
	db := common.GetDB()
	var msg []model.Msg
	db.Raw("SELECT * FROM msgs").Scan(&msg)
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": msg})
}

func DeleteMsg(ctx *gin.Context) {
	sid := ctx.Params.ByName("sid")
	var msg model.Msg
	var user model.User

	db := common.GetDB()
	db.AutoMigrate(model.Msg{})
	db.AutoMigrate(model.User{})

	password := "123456"
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}

	db.Model(&user).Where("sid = ?", sid).Update("password", string(hasedPassword))
	db.Where("sid = ?", sid).Unscoped().Delete(&msg)

	response.Success(ctx, nil, "重置成功")
}

func isExist(db *gorm.DB, sid string) bool {
	var msg model.Msg
	db.Where("sid = ?", sid).First(&msg)
	if msg.ID != 0 {
		return true
	}
	return false
}

func isExist2(db *gorm.DB, sid string) bool {
	var user model.User
	db.Where("sid = ?", sid).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
