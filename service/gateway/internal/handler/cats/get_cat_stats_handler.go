// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package cats

import (
	"net/http"

	"github.com/luyb177/meow-nook/service/gateway/internal/logic/cats"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCatStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := cats.NewGetCatStatsLogic(r.Context(), svcCtx)
		resp, err := l.GetCatStats()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
