// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"net/http"

	"github.com/luyb177/meow-nook/service/gateway/internal/logic/user"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListPointLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewListPointLogsLogic(r.Context(), svcCtx)
		resp, err := l.ListPointLogs()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
