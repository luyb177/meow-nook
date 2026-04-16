// Package adoption contains HTTP handlers for the adoption sub-domain.
package adoption

import (
	"net/http"
	"strconv"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/middleware"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic/adoption"
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

func ListAdoptionsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		var req types.ListAdoptionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := adoption.NewListAdoptionsLogic(r.Context(), ctx).ListAdoptions(claims.UserID, &req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func ApplyAdoptionHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		var req types.ApplyAdoptionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := adoption.NewApplyAdoptionLogic(r.Context(), ctx).ApplyAdoption(claims.UserID, &req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func GetAdoptionHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathID(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		resp, err := adoption.NewGetAdoptionLogic(r.Context(), ctx).GetAdoption(id)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func ReviewAdoptionHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
		var req types.ReviewApplicationReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		if err := adoption.NewReviewAdoptionLogic(r.Context(), ctx).ReviewAdoption(id, claims.UserID, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func ListFollowUpsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathID(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		resp, err := adoption.NewListFollowUpsLogic(r.Context(), ctx).ListFollowUps(id)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func SubmitFollowUpHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
		var req types.SubmitFollowUpReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		if err := adoption.NewSubmitFollowUpLogic(r.Context(), ctx).SubmitFollowUp(id, claims.UserID, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}
