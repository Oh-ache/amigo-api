package baseCode

import (
	"net/http"

	"amigo-api/app/baseCode/api/internal/logic/baseCode"
	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/app/baseCode/rpc/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetItemHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetBaseCodeItemReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := baseCode.NewGetItemLogic(r.Context(), svcCtx)
		resp, err := l.GetItem(&req)
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
