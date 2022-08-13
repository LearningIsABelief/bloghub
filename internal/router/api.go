package router

import (
	"github.com/gin-gonic/gin"
	V1 "gohub/internal/controllers/api/v1"
	login2 "gohub/internal/controllers/api/v1/login"
	register2 "gohub/internal/controllers/api/v1/register"
)

func RegisterRouters(e *gin.Engine) {
	// v1版本的api
	v1 := e.Group("/api/v1")
	// 登录注册
	auth := v1.Group("/auth")
	// 注册
	register := auth.Group("/register")
	{
		register.POST("/using-phone", register2.UsingPhone)
		register.POST("/using-email", register2.UsingEmail)
		register.POST("/phone/exist", register2.PhoneExist)
		register.POST("/email/exist", register2.EmailExist)
		register.POST("/verify-codes/phone", V1.Phone)
		register.POST("/verify-codes/email", V1.Email)
		register.POST("/verify-codes/img", V1.Img)
	}
	// 登录
	login := auth.Group("login")
	{
		login.POST("/using-phone", login2.UsingPhone)
		login.POST("/using-password", login2.UsingPassword)
		login.POST("/password-reset/using-email", login2.PwdResetUsingEmail)
		login.POST("/password-reset/using-phone", login2.PwdResetUsingPhone)
	}
	// 用户管理
	user := v1.Group("/user")
	{
		user.GET("/")
		user.GET("/users")
		user.PUT("/user")
		user.PUT("/user/phone")
		user.PUT("/user/email")
		user.PUT("/user/password")
		// 修改头像
		user.PUT("/user/avatar")
	}
	// 文章管理
	categories := v1.Group("/categories")
	{
		// 分类列表
		categories.GET("/")
		// 创建分类
		categories.POST("/")
		categories.DELETE("/:id")
	}
	// 话题管理
	topics := v1.Group("/topics")
	{
		topics.GET("/")
		topics.POST("/")
		topics.DELETE("/:id")
		topics.PUT("/:id")
		topics.GET("/:id")
	}
	v1.GET("/links")
}
