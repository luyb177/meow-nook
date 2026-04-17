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

func GetAdoptionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdoptionPathIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := adoption.NewGetAdoptionLogic(r.Context(), svcCtx)
		resp, err := l.GetAdoption(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
