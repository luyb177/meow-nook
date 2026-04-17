// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package posts

import (
	"net/http"

	"github.com/luyb177/meow-nook/service/gateway/internal/logic/posts"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UnfavoritePostHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostPathIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := posts.NewUnfavoritePostLogic(r.Context(), svcCtx)
		err := l.UnfavoritePost(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
