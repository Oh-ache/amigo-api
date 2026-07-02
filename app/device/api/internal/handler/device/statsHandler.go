package device

import (
	"net/http"

	"amigo-api/app/device/api/internal/logic/device"
	"amigo-api/app/device/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func StatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := device.NewStatsLogic(r.Context(), svcCtx)
		resp, err := l.Stats()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
