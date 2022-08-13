package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gohub/init/cusRedis"
	"gohub/init/cusZap"
	"gohub/internal/errmsg"
	"gohub/internal/model/request"
	"gohub/internal/pkg"
	"gohub/internal/response"
	"time"
)

// Phone 发送短信验证码
func Phone(c *gin.Context) {
	sendVerifyCodes(c, 1)
}

func Email(c *gin.Context) {
	sendVerifyCodes(c, 2)
}

func Img(c *gin.Context) {
	sendVerifyCodes(c, 3)
}

func sendVerifyCodes(c *gin.Context, phoneOrEmailOrImg int) {
	var param request.Register
	var base64Img string
	var err error
	if err = c.ShouldBindJSON(&param); err != nil {
		response.Response(c, 400, errmsg.BindFailedMsg)
		return
	}
	if len(param.PhoneOrEmail) == 0 {
		response.Response(c, 400, errmsg.PhoneIsEmptyMsg)
		return
	}
	code := pkg.GenValidateCode(6)

	if false {
		if phoneOrEmailOrImg == 1 {
			// 发送短信验证码
			err = pkg.SendSms(param.PhoneOrEmail, code)
			if err != nil {
				cusZap.Error(errmsg.SendSmsFailedMsg, zap.String("err", err.Error()))
				response.Response(c, 500, "", "code", errmsg.SendSmsFailed)
				return
			}
		} else if phoneOrEmailOrImg == 2 {
			// 发送邮箱验证码
			err = pkg.SendEmail(param.PhoneOrEmail, "注册验证码", code)
			if err != nil {
				cusZap.Error(errmsg.SendEmailFailedMsg, zap.String("err", err.Error()))
				response.Response(c, 500, "", "code", errmsg.SendEmailFailed)
				return
			}
		} else {
			codeId := ""
			codeId, base64Img, err = pkg.CreateCode()
			if err != nil {
				cusZap.Error(errmsg.CreateImgCodeFailedMsg, zap.String("err", err.Error()))
				response.Response(c, 500, "", "code", errmsg.CreateImgCodeFailed)
				return
			}
			code = pkg.GetCodeAnswer(codeId)
			fmt.Printf("code:%v\n", code)
		}
	}
	// 将验证码保存到redis中，过期时间为1分钟
	err = cusRedis.Rdb.Set(param.PhoneOrEmail, code, 1*time.Minute).Err()
	if err != nil {
		cusZap.Error(errmsg.SetCodeRedisFailedMsg, zap.String("err", err.Error()))
		response.Response(c, 500, "", "code", errmsg.SetCodeRedisFailed)
		return
	}
	if phoneOrEmailOrImg == 3 {
		response.Response(c, 200, errmsg.SendCodeSuccessMsg, "base64Img", base64Img)
	} else {
		response.Response(c, 200, errmsg.SendCodeSuccessMsg, "verifyCodes", code)
	}
}
