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
			response.Response(c, 401, errmsg.CodeExpiredMsg)
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
	// 3. 判断用户是否存在
	user, userExist, err := repository.UserExist(param.Phone, 1)
	if err != nil {
		response.Response(c, 500, "", "code", errmsg.MySQLQueryFailed)
		return
	}
	if userExist {
		c.Set("user", user)
		accessTokenString, refreshTokenString, errCode := pkg.Accredit(c)
		if errCode != nil {
			response.Response(c, 500, "", "status", "failed", "code", errCode.Error())
			return
		}
		response.Response(c, 200, "登录成功", "status", "success", "accessTokenString", accessTokenString, "refreshTokenString", refreshTokenString)
	} else {
		response.Response(c, 200, errmsg.PhoneDoesNotExistsMsg, "status", "failed")
	}
	// 4. 将验证码缓存设置为过期
	cusRedis.Rdb.Set(param.Phone, correctCode, 1*time.Millisecond)
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
			response.Response(c, 500, "", "status", "", "code", errmsg.MySQLQueryFailed)
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
		c.Set("user", user)
		accessTokenString, refreshTokenString, errCode := pkg.Accredit(c)
		if errCode != nil {
			response.Response(c, 500, "", "status", "failed", "code", errCode.Error())
			return
		}
		response.Response(c, 200, "登录成功", "status", "success", "accessTokenString", accessTokenString, "refreshTokenString", refreshTokenString)
	} else {
		response.Response(c, 200, errmsg.PhoneDoesNotExistsMsg, "status", "failed")
	}
}

func RefreshToken(c *gin.Context) {
	// 1. 从请求头获取tokenString
	refreshTokenString, err := pkg.GetAuthorization(c)
	if err != nil {
		response.Response(c, 401, err.Error())
		return
	}
	// 2. 获取claims
	claims, err := pkg.ParseToken(refreshTokenString)
	if err != nil {
		return
	}
	// 3. 获取refreshTokenString缓存
	RedisRefreshTokenString, err := cusRedis.Rdb.Get(fmt.Sprintf("refreshTokenString:%v", claims.UserID)).Result()
	if err != nil {
		if err == redis.Nil {
			response.Response(c, 401, "refresh token过期")
			return
		}
		response.Response(c, 500, "", "code", errmsg.GetRedisRefreshTokenFailed)
		return
	}
	if RedisRefreshTokenString != refreshTokenString {
		response.Response(c, 401, "refresh token过期")
		return
	}
	// 用户
	user := &model.User{ID: claims.UserID, Name: claims.UserName}
	// 4. 根据用户信息生成accessTokenString和新的refreshTokenString
	c.Set("user", user)
	accessTokenString, refreshTokenString, errCode := pkg.Accredit(c)
	if errCode != nil {
		response.Response(c, 500, "", "code", errCode.Error())
		return
	}
	response.Response(c, 200, "", "accessTokenString", accessTokenString, "refreshTokenString", refreshTokenString)
}
