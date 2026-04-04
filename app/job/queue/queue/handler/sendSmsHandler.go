package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	dataMap, ok := task.Data["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid data format")
	}

	// 将 map 转换为 PushContext
	dataBytes, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	var data message.PushContext
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return err
	}

	sendType, _ := task.Data["send_type"].(string)
	code, _ := task.Data["code"].(string)

	if err := message.PushMessage(&data); err != nil {
		return err
	}

	redisKey := fmt.Sprintf("%s%s:%s", utils.SEND_CODE_KEY, sendType, data.Mobile)
	RedisClient.Set(ctx, redisKey, code, 180*time.Second)

	return nil
}
