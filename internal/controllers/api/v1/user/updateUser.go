package user

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gohub/init/cusRedis"
	"gohub/init/cusZap"
	"gohub/internal/errmsg"
	"gohub/internal/model"
	"gohub/internal/model/request"
	"gohub/internal/repository"
	"gohub/internal/response"
	"golang.org/x/crypto/bcrypt"
	"time"
)

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
			response.Response(c, 401, errmsg.CodeExpiredMsg)
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

func Phone(c *gin.Context) {
	phoneOrEmailReset(c, 1)
}

func Email(c *gin.Context) {
	phoneOrEmailReset(c, 2)
}

type a struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func test(c *gin.Context) {
	w := &a{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
	c.Writer = w
}

func phoneOrEmailReset(c *gin.Context, phoneOrEmail int) {
	var param request.Update
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Response(c, 400, errmsg.BindFailedMsg)
		return
	}
	userID := c.GetUint("userID")
	if err := cusRedis.Rdb.Get(param.NewPhoneOrEmail).Err(); err != nil {
		if err == redis.Nil {
			response.Response(c, 403, "验证码错误")
			return
		}
		response.Response(c, 500, "")
		return
	}
	field := "phone"
	if phoneOrEmail == 2 {
		field = "email"
	}
	err := repository.Update(userID, field, param.NewPhoneOrEmail)
	if err != nil {
		response.Response(c, 500, "")
		return
	}
	response.Response(c, 200, "修改成功")
	cusRedis.Rdb.Set(param.NewPhoneOrEmail, "", time.Millisecond)
}

func UploadImg(c *gin.Context) {
	file, err := c.FormFile("img")
	userID := c.GetUint("userID")
	if err != nil {
		response.Response(c, 500, "上传图片出错")
		return
	}
	fileType := file.Header.Get("Content-Type")
	switch fileType {
	case "image/jpeg":
		fileType = ".jpg"
	case "image/png":
		fileType = ".png"
	default:
		response.Response(c, 403, "文件格式错误")
		return
	}
	avatarPath := fmt.Sprintf("%v/%v/%v%v", viper.GetString("cusPath"), "../img", userID, fileType)
	err = c.SaveUploadedFile(file, avatarPath)
	if err != nil {
		fmt.Println(err)
		response.Response(c, 500, "图片保存失败")
		return
	}
	err = repository.Update(userID, "avatar", avatarPath)
	if err != nil {
		fmt.Println(err)
		response.Response(c, 500, "数据库更新失败")
		return
	}
	response.Response(c, 200, "头像上传成功")
}
