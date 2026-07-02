// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"time"

	"amigo-api/app/job/mqueue/internal/svc"
	"amigo-api/app/job/mqueue/internal/types"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
)

type BoardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBoardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BoardLogic {
	return &BoardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BoardLogic) Board() (resp *types.CommonResp, err error) {
	queues, err := l.svcCtx.Inspector.Queues()
	if err != nil {
		return nil, err
	}

	var total types.QueueStatItem
	items := make([]types.QueueStatItem, 0, len(queues))

	for _, q := range queues {
		pending, _ := l.svcCtx.Inspector.ListPendingTasks(q)
		active, _ := l.svcCtx.Inspector.ListActiveTasks(q)
		scheduled, _ := l.svcCtx.Inspector.ListScheduledTasks(q)
		retry, _ := l.svcCtx.Inspector.ListRetryTasks(q)
		archived, _ := l.svcCtx.Inspector.ListArchivedTasks(q)

		history, _ := l.svcCtx.Inspector.History(q, 1)
		var processed, failed int64
		if len(history) > 0 {
			processed = int64(history[0].Processed)
			failed = int64(history[0].Failed)
		}

		item := types.QueueStatItem{
			Queue:     q,
			Pending:   int64(len(pending)),
			Active:    int64(len(active)),
			Scheduled: int64(len(scheduled)),
			Retry:     int64(len(retry)),
			Archived:  int64(len(archived)),
			Processed: processed,
			Failed:    failed,
		}
		items = append(items, item)
		total.Pending += item.Pending
		total.Active += item.Active
		total.Scheduled += item.Scheduled
		total.Retry += item.Retry
		total.Archived += item.Archived
		total.Processed += item.Processed
		total.Failed += item.Failed
	}

	stat := &types.QueueStatResp{
		Total:     total,
		Queues:    items,
		UpdatedAt: time.Now().Unix(),
	}

	if l.svcCtx.Config.MQueue.ServerName != "" {
		stat.ServerName = l.svcCtx.Config.MQueue.ServerName
	}

	return &types.CommonResp{Code: 0, Data: stat}, nil
}

var _ = asynq.MaxRetry
