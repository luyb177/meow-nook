// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package task

import (
	"net/http"

	"github.com/luyb177/meow-nook/service/gateway/internal/logic/task"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ClaimTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ClaimTaskReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := task.NewClaimTaskLogic(r.Context(), svcCtx)
		resp, err := l.ClaimTask(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
