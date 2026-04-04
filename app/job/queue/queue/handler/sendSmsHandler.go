package handler

import (
	"context"
	"fmt"

	"amigo-api/common/queue"
	"amigo-api/common/utils"
	"amigo-api/common/utils/plug/message"
)

type SendSmsHandler struct{}

func (h *SendSmsHandler) Name() string {
	return "send_sms"
}

func (h *SendSmsHandler) Handle(ctx context.Context, task *queue.Task) error {
	// 从task.Data中获取参数
	// to, _ := task.Data["to"].(string)
	mobile, _ := task.Data["mobile"].(string)
	sendType, _ := task.Data["send_type"].(string)

	pushContext := &message.PushContext{}

	if err := message.PushMessage(pushContext); err != nil {
		return nil
	}

	redisKey := fmt.Sprintf("%s%s:%s", utils.SEND_CODE_KEY, sendType, mobile)
	fmt.Println(redisKey)

	return nil
}
