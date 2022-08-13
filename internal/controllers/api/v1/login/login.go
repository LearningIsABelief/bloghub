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
	"gohub/internal/repository"
	"gohub/internal/response"
	"golang.org/x/crypto/bcrypt"
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
	cusRedis.Rdb.Set(param.Phone, correctCode, -1)
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
	for i := 1; i <= 3 && !userExist; i++ {
		phoneExist := false
		var err error
		user, phoneExist, err = repository.UserExist(param.PhoneOrEmailOrName, i)
		userExist = userExist || phoneExist
		if err != nil {
			response.Response(c, 500, "", "code", errmsg.MySQLQueryFailed)
			return
		}
	}
	// 3. 判断密码是否正确
	if userExist {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.Password)); err != nil {
			fmt.Printf("err:%v\n", err)
			response.Response(c, 200, errmsg.PwdWrongMsg, "status", "failed")
			return
		}
		response.Response(c, 200, errmsg.LoginSuccessMsg, "status", "success")
	} else {
		response.Response(c, 200, errmsg.PhoneDoesNotExistsMsg, "status", "failed")
	}
}

func PwdResetUsingPhone(c *gin.Context) {
	// 1. 参数绑定
	var param request.PwdReset
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Response(c, 400, errmsg.BindFailedMsg)
		return
	}
	// 2. 判断用户是否存在
	var user *model.User
	user, userExist, err := repository.UserExist(param.PhoneOrEmail, 1)
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
	cusRedis.Rdb.Set(param.PhoneOrEmail, correctCode, -1)
	// 4. 更新用户
	user.Password = param.NewPassword
	if userExist {
		repository.UpdateAUser(user)
	}
}

func PwdResetUsingEmail(c *gin.Context) {

}
