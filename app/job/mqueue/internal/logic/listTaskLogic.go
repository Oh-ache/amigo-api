// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"amigo-api/app/job/mqueue/internal/svc"
	"amigo-api/app/job/mqueue/internal/types"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTaskLogic {
	return &ListTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTaskLogic) ListTask(req *types.ListTaskReq) (resp *types.ListTaskResp, err error) {
	page := int(req.Page)
	if page < 1 {
		page = 1
	}
	size := int(req.PageSize)
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	// 决定要查的队列列表
	var queueNames []string
	if req.Queue != "" {
		queueNames = []string{req.Queue}
	} else {
		queues, qerr := l.svcCtx.Inspector.Queues()
		if qerr != nil {
			return nil, qerr
		}
		queueNames = queues
	}

	if len(queueNames) == 0 {
		return &types.ListTaskResp{List: []types.TaskItem{}, Total: 0}, nil
	}

	opts := []asynq.ListOption{
		asynq.PageSize(size),
		asynq.Page(page - 1),
	}
	if req.Queue != "" {
		opts = append(opts, asynq.Queue(req.Queue))
	}

	listFn := func(q string) ([]*asynq.TaskInfo, error) {
		switch req.Status {
		case "active":
			return l.svcCtx.Inspector.ListActiveTasks(q, opts...)
		case "scheduled":
			return l.svcCtx.Inspector.ListScheduledTasks(q, opts...)
		case "retry":
			return l.svcCtx.Inspector.ListRetryTasks(q, opts...)
		case "archived":
			return l.svcCtx.Inspector.ListArchivedTasks(q, opts...)
		case "completed":
			return l.svcCtx.Inspector.ListCompletedTasks(q, opts...)
		default:
			return l.svcCtx.Inspector.ListPendingTasks(q, opts...)
		}
	}

	var allTasks []*asynq.TaskInfo
	for _, q := range queueNames {
		ts, err := listFn(q)
		if err != nil {
			return nil, err
		}
		allTasks = append(allTasks, ts...)
	}

	// 限制总条数
	if len(allTasks) > size {
		allTasks = allTasks[:size]
	}

	items := make([]types.TaskItem, 0, len(allTasks))
	for _, t := range allTasks {
		ti := types.TaskItem{
			TaskID:        t.ID,
			Queue:         t.Queue,
			Type:          t.Type,
			Payload:       string(t.Payload),
			NextProcessAt: t.NextProcessAt.Unix(),
			RetryCount:    int64(t.Retried),
			Status:        req.Status,
		}
		if req.Status == "" {
			ti.Status = "pending"
		}
		items = append(items, ti)
	}

	return &types.ListTaskResp{
		List:  items,
		Total: int64(len(items)),
	}, nil
}
