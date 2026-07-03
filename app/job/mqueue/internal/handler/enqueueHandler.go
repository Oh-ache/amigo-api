// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"amigo-api/app/job/mqueue/internal/logic"
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

		l := logic.NewEnqueueLogic(r.Context(), svcCtx)
		resp, err := l.Enqueue(&req)
		if err != nil {
			result := &types.CommonResp{Code: 1, Msg: err.Error()}
			httpx.OkJsonCtx(r.Context(), w, result)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
