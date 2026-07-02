package event

import (
	"net/http"

	"amigo-api/app/device/api/internal/logic/event"
	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDeviceEventReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := event.NewGetLogic(r.Context(), svcCtx)
		resp, err := l.Get(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
