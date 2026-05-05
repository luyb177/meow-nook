// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package algorithm

import (
	"net/http"

	"github.com/luyb177/meow-nook/service/gateway/internal/logic/algorithm"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 检测猫咪热点区域
func DetectHotspotsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DetectHotspotsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := algorithm.NewDetectHotspotsLogic(r.Context(), svcCtx)
		resp, err := l.DetectHotspots(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
