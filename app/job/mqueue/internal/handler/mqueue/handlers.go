package mqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"amigo-api/common/mqueue"
	"amigo-api/common/utils"
	"amigo-api/common/utils/plug/message"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

// SampleHandler is a sample task handler
type SampleHandler struct{}

func NewSampleHandler() *SampleHandler {
	return &SampleHandler{}
}

func (h *SampleHandler) Name() string {
	return "sample"
}

func (h *SampleHandler) Handle(ctx context.Context, task *mqueue.Task) error {
	fmt.Printf("Processing sample task: %s, data: %v\n", task.ID, task.Data)
	// Simulate some work
	// Add your business logic here
	return nil
}

// EmailHandler handles email sending tasks
type EmailHandler struct{}

func NewEmailHandler() *EmailHandler {
	return &EmailHandler{}
}

func (h *EmailHandler) Name() string {
	return "email"
}

func (h *EmailHandler) Handle(ctx context.Context, task *mqueue.Task) error {
	fmt.Printf("Processing email task: %s, data: %v\n", task.ID, task.Data)
	// Add your email sending logic here
	// e.g., send email using task.Data["to"], task.Data["subject"], etc.
	return nil
}

// NotificationHandler handles notification tasks
type NotificationHandler struct{}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (h *NotificationHandler) Name() string {
	return "notification"
}

func (h *NotificationHandler) Handle(ctx context.Context, task *mqueue.Task) error {
	fmt.Printf("Processing notification task: %s, data: %v\n", task.ID, task.Data)
	// Add your notification logic here
	return nil
}

var RedisClient *redis.Client

// SendSmsHandler handles SMS sending tasks
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

func InitRedis(client *redis.Client) {
	RedisClient = client
}

// TaskHandler handles generic task processing
type TaskHandler struct{}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

func (h *TaskHandler) Name() string {
	return "task"
}

func (h *TaskHandler) Handle(ctx context.Context, task *mqueue.Task) error {
	fmt.Printf("Processing generic task: %s, handler: %s, data: %v\n", task.ID, task.Handler, task.Data)
	return nil
}

