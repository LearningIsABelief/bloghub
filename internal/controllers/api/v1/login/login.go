package login

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
	"time"
)

func UsingPhone(c *gin.Context) {
	// 1. 参数绑定
	var param request.LoginUsingPhone
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Response(c, 400, errmsg.BindFailedMsg)
		return
	}
	// 2. 验证码
	// 2.1 获取验证码缓存
	correctCode, err := cusRedis.Rdb.Get(param.Phone).Result()
	if err != nil {
		if err == redis.Nil {
			response.Response(c, 401, errmsg.PhoneCodeExpiredMsg)
			return
		}
		cusZap.Error(errmsg.GetCodeRedisFailedMsg, zap.String("err", err.Error()))
		response.Response(c, 500, errmsg.LoginFailedMsg, "code", errmsg.GetCodeRedisFailed)
		return
	}
	// 2.2 验证验证码
	if correctCode != param.Code {
		response.Response(c, 403, errmsg.PhoneCodeWrongMsg)
		return
	}
	// 2.3 将验证码缓存设置为过期
	cusRedis.Rdb.Set(param.Phone, correctCode, 1*time.Millisecond)
	// 3. 判断用户是否存在
	_, userExist, err := repository.UserExist(param.Phone, 1)
	if err != nil {
		response.Response(c, 500, "", "code", errmsg.MySQLQueryFailed)
		return
	}
	if userExist {
		response.Response(c, 200, errmsg.LoginSuccessMsg, "status", "success")
	} else {
		response.Response(c, 200, errmsg.PhoneDoesNotExistsMsg, "status", "failed")
	}
}

func UsingPassword(c *gin.Context) {
	// 1. 参数绑定
	var param request.LoginUsingPassword
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Response(c, 400, errmsg.BindFailedMsg)
		return
	}
	// 2. 判断用户是否存在
	var user *model.User
	userExist := false
	isEmail := pkg.VerifyEmailFormat(param.PhoneOrEmailOrName)
	start, step := 1, 2
	if isEmail {
		start, step = 2, 1
	}
	for ; start <= 3 && !userExist; start += step {
		exist := false
		var err error
		user, exist, err = repository.UserExist(param.PhoneOrEmailOrName, start)
		userExist = userExist || exist
		if err != nil {
			response.Response(c, 500, "", "code", errmsg.MySQLQueryFailed)
			return
		}
	}
	// 3. 判断密码是否正确
	if userExist {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.Password)); err != nil {
			fmt.Println(user.Password)
			cusZap.Error("密码比对错误", zap.String("err", err.Error()))
			response.Response(c, 200, errmsg.PwdWrongMsg, "status", "failed")
			return
		}
		response.Response(c, 200, errmsg.LoginSuccessMsg, "status", "success")
	} else {
		response.Response(c, 200, errmsg.PhoneDoesNotExistsMsg, "status", "failed")
	}
}

func PwdResetUsingPhone(c *gin.Context) {
	pwdReset(c, 1)
}

func PwdResetUsingEmail(c *gin.Context) {
	pwdReset(c, 2)
}

func pwdReset(c *gin.Context, phoneOrEmail int) {
	// 1. 参数绑定
	var param request.PwdReset
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Response(c, 400, errmsg.BindFailedMsg)
		return
	}
	// 2. 判断用户是否存在
	var user *model.User
	user, userExist, err := repository.UserExist(param.PhoneOrEmail, phoneOrEmail)
	if err != nil {
		response.Response(c, 500, "", "code", errmsg.MySQLQueryFailed)
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
		response.Response(c, 500, errmsg.LoginFailedMsg, "code", errmsg.GetCodeRedisFailed)
		return
	}
	// 3.2 验证验证码
	if correctCode != param.Code {
		response.Response(c, 403, errmsg.PhoneCodeWrongMsg)
		return
	}
	// 3.3 将验证码缓存设置为过期
	cusRedis.Rdb.Set(param.PhoneOrEmail, correctCode, 1*time.Millisecond)
	// 4. 用户密码加密
	bcryptedPwd, err := bcrypt.GenerateFromPassword([]byte(param.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	//
	user.Password = string(bcryptedPwd)
	if userExist {
		err := repository.UpdateAUser(user)
		if err != nil {
			response.Response(c, 500, "", "code", errmsg.MySQLUpdateFailed)
			return
		}
	}
	response.Response(c, 200, "更新成功", "status", "success")
}
