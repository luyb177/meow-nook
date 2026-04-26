// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package cat

import (
	"net/http"

	"github.com/luyb177/meow-nook/service/gateway/internal/logic/cat"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 志愿者申请创建猫咪档案
func ApplyCreateCatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApplyCreateCatReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := cat.NewApplyCreateCatLogic(r.Context(), svcCtx)
		resp, err := l.ApplyCreateCat(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
