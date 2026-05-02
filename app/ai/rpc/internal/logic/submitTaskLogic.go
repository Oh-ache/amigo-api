package logic

import (
	"context"
	"encoding/json"
	"time"

	"amigo-api/app/ai/model"
	"amigo-api/app/ai/rpc/internal/svc"
	"amigo-api/common/mqueue"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitTaskLogic {
	return &SubmitTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SubmitTaskLogic) SubmitTask(in *pb.SubmitTaskReq) (*pb.SubmitTaskResp, error) {
	now := time.Now().Unix()

	taskTypeStr := taskTypeToString(in.TaskType)

	requestInfo := map[string]interface{}{
		"task_type": taskTypeStr,
		"prompt":    in.Prompt,
		"params":    in.Params,
	}
	requestInfoBytes, _ := json.Marshal(requestInfo)

	task := &model.AiTask{
		UserId:      in.UserId,
		TaskType:    taskTypeStr,
		Prompt:      in.Prompt,
		RequestInfo: string(requestInfoBytes),
		Status:      0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	id, err := l.svcCtx.AiTaskModel.Insert(l.ctx, task)
	if err != nil {
		return nil, err
	}

	mqTask := &mqueue.Task{
		Handler: "ai_task",
		Queue:   "critical",
		Data: map[string]interface{}{
			"id":        id,
			"task_type": taskTypeStr,
			"prompt":    in.Prompt,
			"params":    in.Params,
		},
	}

	_, err = mqueue.GetProducer().Enqueue(l.ctx, mqTask)
	if err != nil {
		logx.Errorf("Enqueue ai_task failed: %v", err)
	}

	return &pb.SubmitTaskResp{
		Id:     id,
		TaskId: "",
	}, nil
}

func taskTypeToString(t pb.TaskType) string {
	switch t {
	case pb.TaskType_TEXT_TO_IMAGE:
		return "text_to_image"
	case pb.TaskType_VIDEO:
		return "video"
	case pb.TaskType_AUDIO:
		return "audio"
	default:
		return "text_to_image"
	}
}