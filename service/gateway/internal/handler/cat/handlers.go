// Package cat contains HTTP handlers for the cat sub-domain.
package cat

import (
	"net/http"
	"strconv"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/middleware"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic/cat"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func pathID(r *http.Request) (int64, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return 0, errorx.New(errorx.CodeBadRequest, "缺少路径参数 id")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errorx.New(errorx.CodeBadRequest, "路径参数 id 格式错误")
	}
	return id, nil
}

func ListCatsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListCatsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := cat.NewListCatsLogic(r.Context(), ctx).ListCats(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func CreateCatHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		var req types.CreateCatReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := cat.NewCreateCatLogic(r.Context(), ctx).CreateCat(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func GetCatHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathID(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		resp, err := cat.NewGetCatLogic(r.Context(), ctx).GetCat(id)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func UpdateCatHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathID(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		var req types.UpdateCatReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		if err := cat.NewUpdateCatLogic(r.Context(), ctx).UpdateCat(id, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func DeleteCatHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathID(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		if err := cat.NewDeleteCatLogic(r.Context(), ctx).DeleteCat(id); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func GetCatStatsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := cat.NewGetCatStatsLogic(r.Context(), ctx).GetCatStats()
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func GetHeatmapHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HeatmapReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := cat.NewGetHeatmapLogic(r.Context(), ctx).GetHeatmap(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func ListRescueRecordsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathID(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		resp, err := cat.NewListRescueRecordsLogic(r.Context(), ctx).ListRescueRecords(id)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func AddRescueRecordHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		id, err := pathID(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		var req types.AddRescueRecordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := cat.NewAddRescueRecordLogic(r.Context(), ctx).AddRescueRecord(id, claims.UserID, &req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func ListHealthRecordsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathID(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		resp, err := cat.NewListHealthRecordsLogic(r.Context(), ctx).ListHealthRecords(id)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func AddHealthRecordHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathID(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		var req types.AddHealthRecordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := cat.NewAddHealthRecordLogic(r.Context(), ctx).AddHealthRecord(id, &req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}
