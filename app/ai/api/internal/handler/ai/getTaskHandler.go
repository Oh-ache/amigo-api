package ai

import (
	"net/http"

	"amigo-api/app/ai/api/internal/logic/ai"
	"amigo-api/app/ai/api/internal/svc"
	"amigo-api/app/ai/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetTaskReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := ai.NewGetTaskLogic(r.Context(), svcCtx)
		resp, err := l.GetTask(&req)
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
