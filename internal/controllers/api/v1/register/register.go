package register

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gohub/init/cusRedis"
	"gohub/init/cusZap"
	"gohub/internal/errmsg"
	"gohub/internal/model"
	"gohub/internal/model/request"
	"gohub/internal/pkg"
	"gohub/internal/repository"
	"gohub/internal/response"
	"golang.org/x/crypto/bcrypt"
)

// UsingPhone 使用手机号注册
func UsingPhone(c *gin.Context) {
	register(c, 1)
}

// UsingEmail 使用邮箱注册
func UsingEmail(c *gin.Context) {
	register(c, 2)
}

func register(c *gin.Context, phoneOrEmail int) {
	var param request.Register
	// 1. 参数绑定
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Response(c, 400, errmsg.BindFailedMsg)
		return
	}
	// 2. 参数合法性检查
	if phoneOrEmail == 2 && !pkg.VerifyEmailFormat(param.PhoneOrEmail) {
		response.Response(c, 403, errmsg.EmailFormatWrongMsg)
		return
	}
	if param.Age <= 0 {
		response.Response(c, 400, errmsg.AgeIllegalMsg)
		return
	}
	// 3. 验证码
	// 3.1 获取验证码缓存
	correctCode, err := cusRedis.Rdb.Get(param.PhoneOrEmail).Result()
	if err != nil {
		if err == redis.Nil {
			response.Response(c, 401, errmsg.PhoneCodeExpiredMsg)
			return
		}
		cusZap.Error(errmsg.GetCodeRedisFailedMsg, zap.String("err", err.Error()))
		response.Response(c, 500, errmsg.RegisterFailedMsg, "code", errmsg.GetCodeRedisFailed)
		return
	}
	// 3.2 验证验证码
	if correctCode != param.Code {
		response.Response(c, 403, errmsg.PhoneCodeWrongMsg)
		return
	}
	// 3.3 将验证码缓存设置为过期
	cusRedis.Rdb.Set(param.PhoneOrEmail, correctCode, -1)
	// 4. 判断手机号/邮箱 和用户名是否已存在
	// 4.1 判断手机号/邮箱是否已存在
	user := model.User{Phone: param.PhoneOrEmail, Name: param.Name, Password: param.Password, Age: param.Age}
	if phoneOrEmail == 1 {
		_, userExist, err := repository.UserExist(user.Phone, phoneOrEmail)
		if err != nil {
			response.Response(c, 500, "", "code", errmsg.MySQLQueryFailed)
			return
		}
		if userExist {
			response.Response(c, 403, errmsg.PhoneAlreadyExistsMsg)
			return
		}
	} else {
		_, userExist, err := repository.UserExist(user.Email, 2)
		if err != nil {
			response.Response(c, 500, "", "code", errmsg.MySQLQueryFailed)
			return
		}
		if userExist {
			response.Response(c, 403, errmsg.EmailAlreadyExistsMsg)
			return
		}
	}
	// 4.2 判断用户名是否已存在
	_, userExist, err := repository.UserExist(user.Name, 3)
	if err != nil {
		response.Response(c, 500, "", "code", errmsg.MySQLQueryFailed)
		return
	}
	if userExist {
		response.Response(c, 403, errmsg.NameAlreadyExistsMsg)
		return
	}
	// 5. 用户密码加密
	encryptedPwd, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
	if err != nil {
		cusZap.Error(errmsg.EncryptedPwdFailedMsg, zap.String("err", err.Error()))
		response.Response(c, 500, errmsg.RegisterFailedMsg, "code", errmsg.EncryptedPwdFailed)
		return
	}
	param.Password = string(encryptedPwd)
	fmt.Printf("register param.Password:%v\n", param.Password)
	// 6. 创建用户
	err = repository.CreateAUser(&user)
	if err != nil {
		response.Response(c, 500, errmsg.RegisterFailedMsg, "code", errmsg.MySQLCreateAUserFailed)
		return
	}
	response.Response(c, 200, errmsg.RegisterSuccessMsg)
	return
}

// PhoneExist 检查手机号是否已注册
func PhoneExist(c *gin.Context) {
	exist(c, 1)
}

func EmailExist(c *gin.Context) {
	exist(c, 2)
}

func exist(c *gin.Context, phoneOrEmail int) {
	var param request.Register
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Response(c, 400, errmsg.BindFailedMsg)
		return
	}
	if len(param.PhoneOrEmail) == 0 {
		response.Response(c, 400, errmsg.PhoneIsEmptyMsg)
		return
	}
	_, exist, err := repository.UserExist(param.PhoneOrEmail, phoneOrEmail)
	if err != nil {
		response.Response(c, 500, "", "code", errmsg.MySQLQueryFailed)
		return
	}
	if exist {
		response.Response(c, 200, "存在", "exist", true)
	} else {
		response.Response(c, 200, "不存在", "exist", false)
	}

}
