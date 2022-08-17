package pkg

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gohub/init/cusZap"
)
import "github.com/go-gomail/gomail"

func SendEmail(to, subject, content string) (err error) {
	from := viper.GetString("email.from")
	password := viper.GetString("email.password")
	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, "zdy") // 发件人
	m.SetHeader("To",                       // 收件人
		m.FormatAddress(to, "Gohub"),
	)
	m.SetHeader("Subject", subject)                                // 主题
	m.SetBody("text/html", content)                                // 正文
	d := gomail.NewPlainDialer("smtp.qq.com", 465, from, password) // 发送邮件服务器、端口、发件人账号、发件人密码
	if err = d.DialAndSend(m); err != nil {
		cusZap.Error("邮件发送失败", zap.String("err", err.Error()))
		return
	}
	return
}
