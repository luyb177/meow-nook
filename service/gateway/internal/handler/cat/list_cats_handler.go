// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package cat

import (
	"net/http"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/httpresp"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic/cat"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取小猫列表（支持筛选、分页、附近查找）
func ListCatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListCatsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresp.JsonBaseResponseCtx(r.Context(), w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}

		l := cat.NewListCatsLogic(r.Context(), svcCtx)
		resp, err := l.ListCats(&req)
		if err != nil {
			httpresp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			httpresp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
