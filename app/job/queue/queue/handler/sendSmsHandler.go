package handler

import (
	"context"
	"fmt"
	"time"

	"amigo-api/app/job/queue/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendSmsHandler struct{}

func (h *SendSmsHandler) Name() string {
	return "send_sms"
}

func (h *SendSmsHandler) Handle(ctx context.Context, task *types.Task) error {
	// 从task.Data中获取参数
	to, _ := task.Data["to"].(string)
	subject, _ := task.Data["subject"].(string)
	body, _ := task.Data["body"].(string)

	// 发送逻辑
	logx.Infof("Sending email to %s: %s", to, subject)
	fmt.Println(body)

	// 模拟发送
	time.Sleep(100 * time.Millisecond)

	return nil
}

