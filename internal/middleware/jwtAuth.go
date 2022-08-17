package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gohub/init/cusRedis"
	"gohub/internal/errmsg"
	"gohub/internal/pkg"
	"gohub/internal/response"
)

func Authentication() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 1. 获取accessTokenString
		accessTokenString, err := pkg.GetAuthorization(c)
		if err != nil {
			response.Response(c, 401, err.Error())
			c.Abort()
			return
		}
		// 2. 获取claims
		claims, err := pkg.ParseToken(accessTokenString)
		if err != nil {
			response.Response(c, 401, "登录已过期")
			c.Abort()
			return
		}
		if claims == nil {
			response.Response(c, 401, "无效token")
			c.Abort()
			return
		}
		// 3. 获取在线用户缓存
		redisAccessTokenString, err := cusRedis.Rdb.Get(fmt.Sprintf("onlineUser:%v", claims.UserID)).Result()
		if err != nil {
			if err == redis.Nil {
				response.Response(c, 401, "登录已过期")
				c.Abort()
				return
			}
			response.Response(c, 500, "", "code", errmsg.GetRedisOnlineUserFailed)
			c.Abort()
			return
		}
		if redisAccessTokenString != accessTokenString {
			response.Response(c, 401, "登录已过期")
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("userName", claims.UserName)
		c.Next()
	}
}
