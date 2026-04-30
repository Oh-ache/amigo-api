package mqueue

import (
	"context"
	"encoding/json"
	"time"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
	StatusDeadLetter TaskStatus = "dead_letter"
)

// Priority represents task priority levels
type Priority int

const (
	PriorityLow    Priority = 1
	PriorityNormal Priority = 5
	PriorityHigh   Priority = 10
	PriorityUrgent Priority = 100
)

// Task represents a message queue task
type Task struct {
	ID         string                 `json:"id"`          // Task ID
	Queue      string                 `json:"queue"`       // Queue name
	Data       map[string]interface{} `json:"data"`        // Task data
	Priority   Priority               `json:"priority"`    // Priority
	MaxRetry   int                    `json:"max_retry"`   // Max retry count
	RetryCount int                    `json:"retry_count"` // Current retry count
	CreatedAt  int64                  `json:"created_at"`  // Creation timestamp
	Delay      int64                  `json:"delay"`       // Delay in seconds
	Timeout    int64                  `json:"timeout"`     // Timeout in seconds
	Handler    string                 `json:"handler"`     // Handler name
	TaskType   string                 `json:"task_type"`   // Task type for asynq
}

// TaskResult represents the result of a task execution
type TaskResult struct {
	TaskID     string      `json:"task_id"`
	Status     TaskStatus  `json:"status"`
	Error      string      `json:"error,omitempty"`
	Result     interface{} `json:"result,omitempty"`
	Duration   int64       `json:"duration"` // Execution duration in milliseconds
	FinishedAt int64       `json:"finished_at"`
}

// Handler is the interface for task handlers
type Handler interface {
	Handle(ctx context.Context, task *Task) error
	Name() string
}

// HandlerFunc is a function type that implements Handler
type HandlerFunc func(ctx context.Context, task *Task) error

func (f HandlerFunc) Handle(ctx context.Context, task *Task) error {
	return f(ctx, task)
}

func (f HandlerFunc) Name() string {
	return "anonymous"
}

// ConvertTaskToPayload converts a Task to asynq Payload
func ConvertTaskToPayload(task *Task) (string, error) {
	data, err := json.Marshal(task)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ParsePayloadToTask parses asynq payload back to Task
func ParsePayloadToTask(payload []byte) (*Task, error) {
	var task Task
	if err := json.Unmarshal(payload, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

// GetDefaultQueueConfig returns default queue configuration
func GetDefaultQueueConfig() *QueueConfig {
	return &QueueConfig{
		RedisOpt:      nil,
		Queues:        map[string]int{"default": 6, "critical": 3, "low": 1},
		Concurrency:   10,
		StrictMode:    false,
		SyncTimeout:   10 * time.Second,
		RetryDelay:    1 * time.Minute,
		MaxRetry:      3,
		Timeout:       30 * time.Minute,
		DeadQueueName: "dead",
		ServerName:    "amigo-mqueue",
	}
}

