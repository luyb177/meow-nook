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

// 管理员审核通过
func ApproveCreateCatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApproveCreateCatReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := cat.NewApproveCreateCatLogic(r.Context(), svcCtx)
		resp, err := l.ApproveCreateCat(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
