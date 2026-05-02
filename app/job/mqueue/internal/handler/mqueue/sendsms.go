package mqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"amigo-api/common/mqueue"
	"amigo-api/common/utils"
	"amigo-api/common/utils/plug/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendSmsHandler struct{}

func NewSendSmsHandler() *SendSmsHandler {
	return &SendSmsHandler{}
}

func (h *SendSmsHandler) Name() string {
	return "send_sms"
}

func (h *SendSmsHandler) Handle(ctx context.Context, task *mqueue.Task) error {
	dataMap, ok := task.Data["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid data format")
	}

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

	logx.Infof("SendSmsHandler: SMS sent to %s, code: %s", data.Mobile, code)
	return nil
}
