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
	data, _ := task.Data["data"].(message.PushContext)
	sendType, _ := task.Data["send_type"].(string)
	code, _ := task.Data["code"].(string)

	if err := message.PushMessage(&data); err != nil {
		return nil
	}

	redisKey := fmt.Sprintf("%s%s:%s", utils.SEND_CODE_KEY, sendType, data.Mobile)
	fmt.Println(redisKey)
	RedisClient.Set(ctx, redisKey, code, 180)

	return nil
}
