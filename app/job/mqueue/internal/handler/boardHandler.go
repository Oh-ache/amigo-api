// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"amigo-api/app/job/mqueue/internal/logic"
	"amigo-api/app/job/mqueue/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func BoardHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewBoardLogic(r.Context(), svcCtx)
		resp, err := l.Board()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
