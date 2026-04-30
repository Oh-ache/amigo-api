package handler

import (
	"net/http"

	"amigo-api/app/job/mqueue/internal/svc"
	"amigo-api/app/job/mqueue/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func EnqueueHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EnqueueReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		taskID, err := svcCtx.EnqueueTask(r.Context(), req.Handler, req.Data)
		result := &types.CommonResp{}
		if err != nil {
			result.Code = 1
			result.Msg = err.Error()
			result.Data = &types.EmptyResp{}
		} else {
			result.Code = 0
			result.Msg = "success"
			result.Data = &types.TaskInfoResp{TaskID: taskID}
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func EnqueueDelayedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EnqueueReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		taskID, err := svcCtx.EnqueueDelayedTask(r.Context(), req.Handler, req.Data, req.Delay)
		result := &types.CommonResp{}
		if err != nil {
			result.Code = 1
			result.Msg = err.Error()
			result.Data = &types.EmptyResp{}
		} else {
			result.Code = 0
			result.Msg = "success"
			result.Data = &types.TaskInfoResp{TaskID: taskID}
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func TaskResultHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TaskResultReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		taskResult, err := svcCtx.Consumer.GetTaskResult(r.Context(), req.TaskID)
		result := &types.CommonResp{}
		if err != nil {
			result.Code = 1
			result.Msg = err.Error()
			result.Data = &types.EmptyResp{}
		} else {
			result.Code = 0
			result.Msg = "success"
			result.Data = &types.TaskResultResp{
				TaskID:     taskResult.TaskID,
				Status:     string(taskResult.Status),
				Error:      taskResult.Error,
				Duration:   taskResult.Duration,
				FinishedAt: taskResult.FinishedAt,
			}
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}
