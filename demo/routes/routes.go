package routes

import (
	"demo/controller"
	"demo/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/admin/register", controller.ARegister)

	r.POST("/api/auth/login", controller.Login)
	r.POST("/api/admin/login", controller.ALogin)

	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	r.GET("/api/admin/info", middleware.AdminMiddleware(), controller.AdminInfo)

	r.PUT("/api/admin/:id", middleware.AdminMiddleware(), controller.AUpdate)
	r.DELETE("/api/admin/:id", middleware.AdminMiddleware(), controller.ADelete)
	r.GET("/api/admin/show", middleware.AdminMiddleware(), controller.AShow)
	r.GET("/api/admin/search", middleware.AdminMiddleware(), controller.ASearch)

	r.PUT("/api/auth/:id", middleware.AuthMiddleware(), controller.PassUpdate)
	r.POST("/api/auth/:id", controller.InfoUpdate)
	r.DELETE("/api/auth/:id", middleware.AdminMiddleware(), controller.Delete)

	r.POST("/api/message", controller.ForgetPass)
	r.GET("/api/message/show", middleware.AdminMiddleware(), controller.ShowMsg)
	r.DELETE("/api/message/:sid", middleware.AdminMiddleware(), controller.DeleteMsg)

	r.POST("/api/class/add", middleware.AdminMiddleware(), controller.AddClass)
	r.GET("/api/class/show", middleware.AdminMiddleware(), controller.ShowClass)
	r.DELETE("/api/class/:name", middleware.AdminMiddleware(), controller.DeleteClass)
	r.POST("/api/class/:name", middleware.AdminMiddleware(), controller.UpdateClass)

	return r
}
