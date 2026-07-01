package admin

import (
	"net/http"

	"amigo-api/app/user/api/internal/logic/admin"
	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DomainListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewDomainListLogic(r.Context(), svcCtx)
		resp, err := l.DomainList()
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
