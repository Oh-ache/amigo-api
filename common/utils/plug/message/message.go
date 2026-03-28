package message

import (
	"errors"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type PushContext struct {
	Platform string
	Mobile   string
	Email    string
	Content  string

	TmplateCode string
	SignName    string

	AliAccessKeyId     string
	AliAccessKeySecret string
}

type PaymentStrategy interface {
	Push(ctx *PushContext) error
}

func NewPush(ctx *PushContext, strategy PaymentStrategy) *Push {
	return &Push{
		Content:  ctx,
		Strategy: strategy,
	}
}

type Push struct {
	Content  *PushContext
	Strategy PaymentStrategy
}

func (p *Push) Push(ctx *PushContext) error {
	return p.Strategy.Push(ctx)
}

// 定义发送方式
var sendMethod = make(map[string]PaymentStrategy)

func init() {
	sendMethod["ali_sms"] = &AliyunSms{}
}

func PushMessage(ctx *PushContext) error {
	platform := ctx.Platform

	strategy, ok := sendMethod[platform]

	if !ok {
		return errors.New("invalid platform")
	}

	payment := NewPush(ctx, strategy)
	return payment.Push(ctx)
}

// -------------------------------------阿里云短信
type AliyunSms struct {
}

func (a *AliyunSms) Push(ctx *PushContext) error {
	config := &openapi.Config{
		AccessKeyId:     tea.String(ctx.AliAccessKeyId),
		AccessKeySecret: tea.String(ctx.AliAccessKeySecret),
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	client, _err := dysmsapi20170525.NewClient(config)

	if _err != nil {
		return errors.New("初始化短信客户端失败")
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(ctx.Mobile),
		SignName:      tea.String(ctx.SignName),
		TemplateCode:  tea.String(ctx.TmplateCode),
		TemplateParam: tea.String(ctx.Content),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_, _err = client.SendSmsWithOptions(sendSmsRequest, runtime)
		if _err != nil {
			return _err
		}

		return nil
	}()

	if tryErr != nil {
		return errors.New("短信发送失败")
		// var error = &tea.SDKError{}
		// if _t, ok := tryErr.(*tea.SDKError); ok {
		//     error = _t
		// } else {
		//     error.Message = tea.String(tryErr.Error())
		// }
		// // 此处仅做打印展示，请谨慎对待异常处理，在工程项目中切勿直接忽略异常。
		// // 错误 message
		// fmt.Println(tea.StringValue(error.Message))
		// // 诊断地址
		// var data interface{}
		// d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		// d.Decode(&data)
		// if m, ok := data.(map[string]interface{}); ok {
		//     recommend, _ := m["Recommend"]
		//     fmt.Println(recommend)
		// }
		// _, _err = util.AssertAsString(error.Message)
		// if _err != nil {
		//     return _err
		// }
	}
	return nil
}
