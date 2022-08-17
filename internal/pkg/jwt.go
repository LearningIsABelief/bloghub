package pkg

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gohub/init/cusRedis"
	"gohub/init/cusZap"
	"gohub/internal/errmsg"
	"gohub/internal/model"
	"time"
)

type cusClaims struct {
	UserID   uint
	UserName string
	jwt.RegisteredClaims
}

var secret = viper.GetString("jwt.secret")

func GenToken(user *model.User) (accessTokenString string, accessTokenErr error, refreshTokenString string, refreshTokenErr error) {
	accessClaims := &cusClaims{
		user.ID,
		user.Name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(viper.GetInt("jwt.accessExpireAt")) * time.Minute)),
			Issuer:    viper.GetString("jwt.issuer"),
		},
	}
	refreshClaims := &cusClaims{
		user.ID,
		user.Name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(viper.GetInt("jwt.refreshExpireAt")) * time.Minute)),
			Issuer:    viper.GetString("jwt.issuer"),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	accessTokenString, accessTokenErr = accessToken.SignedString([]byte(secret))
	refreshTokenString, refreshTokenErr = refreshToken.SignedString([]byte(secret))
	return
}

func ParseToken(tokenString string) (*cusClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &cusClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, accessOk := token.Claims.(*cusClaims)
	if accessOk && token.Valid {
		return claims, nil
	}
	return nil, nil
}

func Accredit(c *gin.Context) (accessTokenString, refreshTokenString string, errCode error) {
	var user *model.User
	// 1. 从上下文中获取user
	get, _ := c.Get("user")
	user = get.(*model.User)
	// 2. 获取accessTokenString和refreshTokenString
	accessTokenString, accessTokenErr, refreshTokenString, refreshTokenErr := GenToken(user)
	if accessTokenErr != nil && refreshTokenErr != nil {
		cusZap.Error(errmsg.GenTokenFailedMsg, zap.String("accessTokenErr", accessTokenErr.Error()), zap.String("refreshTokenErr", refreshTokenErr.Error()))
		return "", "", errors.New(string(errmsg.GenTokenFailed))
	}
	// 3. 将accessTokenString保存到redis中
	// onlineUser:用户id作为key，accessTokenString实现单点登录
	err := cusRedis.Rdb.Set(fmt.Sprintf("onlineUser:%v", user.ID), accessTokenString, time.Duration(viper.GetInt("jwt.accessExpireAt"))*time.Minute).Err()
	if err != nil {
		cusZap.Error(errmsg.SetRedisOnlineUserFailedMsg, zap.String("err", err.Error()))
		return "", "", errors.New(string(errmsg.SetRedisOnlineUserFailed))
	}
	// 4. 将refreshTokenString保存到redis中
	err = cusRedis.Rdb.Set(fmt.Sprintf("refreshTokenString:%v", user.ID), refreshTokenString, time.Duration(viper.GetInt("jwt.refreshExpireAt"))*time.Minute).Err()
	if err != nil {
		cusZap.Error(errmsg.SetRedisOnlineUserFailedMsg, zap.String("err", err.Error()))
		return "", "", errors.New(string(errmsg.SetRedisRefreshTokenFailed))
	}
	return
}
