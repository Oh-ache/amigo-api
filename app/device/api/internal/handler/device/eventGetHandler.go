// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package device

import (
	"net/http"

	"amigo-api/app/device/api/internal/logic/device"
	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func EventGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDeviceEventReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := device.NewEventGetLogic(r.Context(), svcCtx)
		resp, err := l.EventGet(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
