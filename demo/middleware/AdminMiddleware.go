package middleware

import (
	"demo/common"
	"demo/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		//validate token format 如果为空或者不以Bearer开头
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}
		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		//验证通过后获取claims中的user.Sid或admin.Email
		userId := claims.UserId
		DB := common.GetDB()
		var admin model.Admin
		DB.First(&admin, userId)

		//用户不存在
		if admin.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		//用户存在 将信息写入上下文
		ctx.Set("admin", admin)

		ctx.Next()
	}
}
