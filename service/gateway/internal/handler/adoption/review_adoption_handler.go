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

func ReviewAdoptionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReviewApplicationReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := adoption.NewReviewAdoptionLogic(r.Context(), svcCtx)
		err := l.ReviewAdoption(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
