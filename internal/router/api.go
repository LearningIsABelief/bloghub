package router

import (
	"github.com/gin-gonic/gin"
	loginHandler "gohub/internal/controllers/api/v1/login"
	registerHandler "gohub/internal/controllers/api/v1/register"
	userHandler "gohub/internal/controllers/api/v1/user"
	verifyCodesHandler "gohub/internal/controllers/api/v1/verifyCodes"
	"gohub/internal/middleware"
)

func RegisterRouters(e *gin.Engine) {
	// v1版本的api
	v1 := e.Group("/api/v1")
	// 登录注册
	auth := v1.Group("/auth")
	// 注册
	register := auth.Group("/register")
	{
		register.POST("/using-phone", registerHandler.UsingPhone)
		register.POST("/using-email", registerHandler.UsingEmail)
		register.POST("/phone/exist", registerHandler.PhoneExist)
		register.POST("/email/exist", registerHandler.EmailExist)
	}
	// 登录
	login := auth.Group("login")
	{
		login.GET("/refresh-token", loginHandler.RefreshToken)
		login.POST("/using-phone", loginHandler.UsingPhone)
		login.POST("/using-password", loginHandler.UsingPassword)
	}
	// 获取验证码
	getVerifyCode := auth.Group("/verify-codes")
	{
		getVerifyCode.POST("/phone", verifyCodesHandler.Phone)
		getVerifyCode.POST("/email", verifyCodesHandler.Email)
		getVerifyCode.POST("/img", verifyCodesHandler.Img)
	}
	// 用户管理
	needAuth := v1.Group("", middleware.Authentication())
	user := needAuth.Group("/user")
	{
		user.GET("/", userHandler.User)
		user.GET("/all", userHandler.All)
		user.PUT("/phone", userHandler.Phone)
		user.PUT("/email", userHandler.Email)
		v1.PUT("/user/password/using-email", userHandler.PwdResetUsingEmail)
		v1.PUT("/user/password/using-phone", userHandler.PwdResetUsingPhone)
		// 修改头像
		user.PUT("/avatar", userHandler.UploadImg)
	}
	// 文章管理
	categories := needAuth.Group("/categories")
	{
		// 分类列表
		categories.GET("/")
		// 创建分类
		categories.POST("/")
		categories.DELETE("/:id")
	}
	// 话题管理
	topics := needAuth.Group("/topics")
	{
		topics.GET("/")
		topics.POST("/")
		topics.DELETE("/:id")
		topics.PUT("/:id")
		topics.GET("/:id")
	}
	needAuth.GET("/links")
}
