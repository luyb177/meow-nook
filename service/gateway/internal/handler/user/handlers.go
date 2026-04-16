// Package user contains HTTP handlers for the user sub-domain.
package user

import (
	"net/http"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/middleware"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic/user"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RegisterHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := user.NewRegisterLogic(r.Context(), ctx).Register(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func LoginHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := user.NewLoginLogic(r.Context(), ctx).Login(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func LogoutHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		if err := user.NewLogoutLogic(r.Context(), ctx).Logout(claims.UserID); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func GetUserInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		resp, err := user.NewGetUserInfoLogic(r.Context(), ctx).GetUserInfo(claims.UserID)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func UpdateUserInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		var req types.UpdateUserInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		if err := user.NewUpdateUserInfoLogic(r.Context(), ctx).UpdateUserInfo(claims.UserID, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func ChangePasswordHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		var req types.ChangePasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		if err := user.NewChangePasswordLogic(r.Context(), ctx).ChangePassword(claims.UserID, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func GetPointsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		resp, err := user.NewGetPointsLogic(r.Context(), ctx).GetPoints(claims.UserID)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func ListPointLogsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		var req types.PageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := user.NewListPointLogsLogic(r.Context(), ctx).ListPointLogs(claims.UserID, &req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func GetNotificationSettingsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		resp, err := user.NewGetNotificationSettingsLogic(r.Context(), ctx).GetNotificationSettings(claims.UserID)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func UpdateNotificationSettingsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		var req types.NotificationSettings
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		if err := user.NewUpdateNotificationSettingsLogic(r.Context(), ctx).UpdateNotificationSettings(claims.UserID, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func SubmitFeedbackHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		var req types.FeedbackReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := user.NewSubmitFeedbackLogic(r.Context(), ctx).SubmitFeedback(claims.UserID, &req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}
