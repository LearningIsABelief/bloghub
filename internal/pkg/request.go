package pkg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gohub/internal/errmsg"
	"strings"
)

func GetAuthorization(c *gin.Context) (string, error) {
	// 1. 从请求头中获取授权authorization
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		return "", fmt.Errorf(errmsg.AuthIsEmptyMsg)
	}
	// 2. 获取refreshTokenString
	parts := strings.SplitN(authorization, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", fmt.Errorf("请求头中的auth格式错误")
	}
	tokenString := parts[1]
	return tokenString, nil
}
