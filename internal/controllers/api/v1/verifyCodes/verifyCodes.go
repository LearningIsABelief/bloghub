package verifyCodes

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
	send(c, 1)
}

func Email(c *gin.Context) {
	send(c, 2)
}

func Img(c *gin.Context) {
	send(c, 3)
}

func send(c *gin.Context, phoneOrEmailOrImg int) {
	var param request.VerifyCodes
	var base64Img string
	var err error
	if phoneOrEmailOrImg != 3 && c.ShouldBindJSON(&param) != nil {
		response.Response(c, 400, errmsg.BindFailedMsg)
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
		}
	}
	// 将验证码保存到redis中，过期时间为1分钟
	err = cusRedis.Rdb.Set(param.PhoneOrEmail, code, time.Duration(viper.GetInt("verifyCode.expireAt"))*time.Minute).Err()
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
