package oss

import (
	"net/http"

	"amigo-api/app/sdk/api/internal/logic/sdk/oss"
	"amigo-api/app/sdk/api/internal/svc"
	"amigo-api/app/sdk/api/internal/types"
	"amigo/app/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := oss.NewUploadTokenLogic(r.Context(), svcCtx)
		resp, err := l.UploadToken(&req)
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
