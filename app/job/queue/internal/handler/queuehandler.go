package handler

import (
	"net/http"

	"amigo-api/app/job/queue/internal/logic"
	"amigo-api/app/job/queue/internal/svc"
	"amigo-api/app/job/queue/internal/types"
	"amigo/app/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func QueueHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewQueueLogic(r.Context(), svcCtx)
		resp, err := l.Queue(&req)
		result := &types.CommonResp{}
		if err != nil {
			result.Code = 1
			result.Msg = err.Error()
			result.Data = &types.EmptyResp{}
		} else {
			result.Data = resp
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}
