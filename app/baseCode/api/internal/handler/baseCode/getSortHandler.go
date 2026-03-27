package baseCode

import (
	"net/http"

	"amigo-api/app/baseCode/api/internal/logic/baseCode"
	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/app/baseCode/rpc/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetSortHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetBaseCodeSortReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := baseCode.NewGetSortLogic(r.Context(), svcCtx)
		resp, err := l.GetSort(&req)
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
