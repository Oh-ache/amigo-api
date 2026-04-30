package mqueue

import (
	"context"
	"fmt"

	"amigo-api/common/mqueue"
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
	// Add your generic task logic here
	return nil
}

