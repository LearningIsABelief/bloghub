// This file is auto-generated, don't edit it. Thanks.
package pkg

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"go.uber.org/zap"
	"gohub/init/cusZap"
	"os"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func _main(args []*string, phone, code string) (_err error) {
	client, _err := CreateClient(tea.String("LTAI5t63RPHg2MmSMuT1mtSf"), tea.String("XQ8jyc9azcJBnWSOCGh4kRXTa20CUx"))
	if _err != nil {
		return _err
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("阿里云短信测试"),
		TemplateCode:  tea.String("SMS_154950909"),
		PhoneNumbers:  tea.String(phone),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%v\"}", code)),
	}
	runtime := &util.RuntimeOptions{}
	resp, _err := client.SendSmsWithOptions(sendSmsRequest, runtime)
	if _err != nil {
		return _err
	}
	if *(resp.Body.Code) == "OK" {
		fmt.Printf("发送状态：%v\n", *resp.Body.Code)
	} else {
		_err = fmt.Errorf(*util.ToJSONString(tea.ToMap(resp)))
	}
	return _err
}

func SendSms(phone, code string) error {
	err := _main(tea.StringSlice(os.Args[1:]), phone, code)
	if err != nil {
		cusZap.Logger.Error("短信发送失败", zap.String("err", err.Error()))
		return err
	}
	return nil
}
