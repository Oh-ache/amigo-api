package handler

import (
	"net/http"

	"amigo-api/app/job/queue/internal/logic"
	"amigo-api/app/job/queue/internal/svc"
	"amigo-api/app/job/queue/internal/types"

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
		resp, _ := l.Queue(&req)
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
