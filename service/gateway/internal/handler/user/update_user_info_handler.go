package user

import (
	"net/http"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/httpresp"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic/user"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateUserInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateUserInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresp.JsonBaseResponseCtx(r.Context(), w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}

		l := user.NewUpdateUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.UpdateUserInfo(&req)
		if err != nil {
			httpresp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			httpresp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}

