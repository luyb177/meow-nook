// Package post contains HTTP handlers for the post sub-domain.
package post

import (
	"net/http"
	"strconv"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/middleware"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic/post"
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

func pathCommentID(r *http.Request) (int64, error) {
	cidStr := r.PathValue("cid")
	if cidStr == "" {
		return 0, errorx.New(errorx.CodeBadRequest, "缺少路径参数 cid")
	}
	cid, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil {
		return 0, errorx.New(errorx.CodeBadRequest, "路径参数 cid 格式错误")
	}
	return cid, nil
}

func ListPostsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListPostsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := post.NewListPostsLogic(r.Context(), ctx).ListPosts(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func CreatePostHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		var req types.CreatePostReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := post.NewCreatePostLogic(r.Context(), ctx).CreatePost(claims.UserID, &req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func GetPostHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathID(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		resp, err := post.NewGetPostLogic(r.Context(), ctx).GetPost(id)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func DeletePostHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
		if err := post.NewDeletePostLogic(r.Context(), ctx).DeletePost(id, claims.UserID); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func LikePostHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
		if err := post.NewLikePostLogic(r.Context(), ctx).LikePost(id, claims.UserID); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func UnlikePostHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
		if err := post.NewUnlikePostLogic(r.Context(), ctx).UnlikePost(id, claims.UserID); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func FavoritePostHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
		if err := post.NewFavoritePostLogic(r.Context(), ctx).FavoritePost(id, claims.UserID); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func UnfavoritePostHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
		if err := post.NewUnfavoritePostLogic(r.Context(), ctx).UnfavoritePost(id, claims.UserID); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func ListCommentsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathID(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		var req types.PageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := post.NewListCommentsLogic(r.Context(), ctx).ListComments(id, &req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func AddCommentHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
		var req types.AddCommentReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		if err := post.NewAddCommentLogic(r.Context(), ctx).AddComment(id, claims.UserID, req.Content); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func DeleteCommentHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		cid, err := pathCommentID(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		if err := post.NewDeleteCommentLogic(r.Context(), ctx).DeleteComment(cid, claims.UserID); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func ReportPostHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
		var req types.ReportPostReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		if err := post.NewReportPostLogic(r.Context(), ctx).ReportPost(id, claims.UserID, req.Reason); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}
