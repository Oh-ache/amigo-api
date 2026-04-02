package user

import (
	"net/http"

	"amigo-api/app/user/api/internal/logic/user"
	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LoginMiniappHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserLoginMiniappReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewLoginMiniappLogic(r.Context(), svcCtx)
		resp, err := l.LoginMiniapp(&req)
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
