package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gohub/internal/repository"
	"gohub/internal/response"
)

func User(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := repository.User(userID)
	if err != nil {
		response.Response(c, 500, "")
		return
	}
	response.Response(c, 200, "获取成功", "user", user)
}

func All(c *gin.Context) {
	users, err := repository.All()
	if err != nil {
		response.Response(c, 500, "")
		return
	}
	response.Response(c, 200, fmt.Sprintf("共%v条用户信息", len(users)), "users", users)
}
