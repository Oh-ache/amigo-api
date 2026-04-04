package queue

import "context"

// 任务状态
type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
	StatusDeadLetter TaskStatus = "dead_letter"
)

// 任务优先级
type Priority int

const (
	PriorityLow    Priority = 1
	PriorityNormal Priority = 5
	PriorityHigh   Priority = 10
	PriorityUrgent Priority = 100
)

// 任务定义
type Task struct {
	ID         string                 `json:"id"`          // 任务ID
	Queue      string                 `json:"queue"`       // 队列名称
	Data       map[string]interface{} `json:"data"`        // 任务数据
	Priority   Priority               `json:"priority"`    // 优先级
	MaxRetry   int                    `json:"max_retry"`   // 最大重试次数
	RetryCount int                    `json:"retry_count"` // 当前重试次数
	CreatedAt  int64                  `json:"created_at"`  // 创建时间戳
	Delay      int64                  `json:"delay"`       // 延迟执行秒数
	Timeout    int64                  `json:"timeout"`     // 超时时间(秒)
	Handler    string                 `json:"handler"`     // 处理器名称
}

// 任务结果
type TaskResult struct {
	TaskID     string      `json:"task_id"`
	Status     TaskStatus  `json:"status"`
	Error      string      `json:"error,omitempty"`
	Result     interface{} `json:"result,omitempty"`
	Duration   int64       `json:"duration"` // 执行耗时(毫秒)
	FinishedAt int64       `json:"finished_at"`
}

// 处理器接口
type Handler interface {
	Handle(ctx context.Context, task *Task) error
	Name() string
}

// 处理器函数
type HandlerFunc func(ctx context.Context, task *Task) error

func (f HandlerFunc) Handle(ctx context.Context, task *Task) error {
	return f(ctx, task)
}

func (f HandlerFunc) Name() string {
	return "anonymous"
}
