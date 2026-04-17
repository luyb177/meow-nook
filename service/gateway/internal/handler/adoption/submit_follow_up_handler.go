// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package adoption

import (
	"net/http"

	"github.com/luyb177/meow-nook/service/gateway/internal/logic/adoption"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SubmitFollowUpHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SubmitFollowUpReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := adoption.NewSubmitFollowUpLogic(r.Context(), svcCtx)
		err := l.SubmitFollowUp(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
