// Package task contains HTTP handlers for the task sub-domain.
package task

import (
	"net/http"
	"strconv"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/middleware"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic/task"
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

func ListTasksHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListTasksReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := task.NewListTasksLogic(r.Context(), ctx).ListTasks(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func CreateTaskHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.ClaimsFromContext(r.Context())
		if claims == nil {
			httpx.Error(w, errorx.ErrUnauthorized)
			return
		}
		var req types.CreateTaskReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		resp, err := task.NewCreateTaskLogic(r.Context(), ctx).CreateTask(claims.UserID, &req)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func GetTaskHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathID(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		resp, err := task.NewGetTaskLogic(r.Context(), ctx).GetTask(id)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func UpdateTaskHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
		var req types.UpdateTaskReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		if err := task.NewUpdateTaskLogic(r.Context(), ctx).UpdateTask(id, claims.UserID, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func ClaimTaskHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
		if err := task.NewClaimTaskLogic(r.Context(), ctx).ClaimTask(id, claims.UserID); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func AbandonTaskHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
		var req types.AbandonTaskReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		if err := task.NewAbandonTaskLogic(r.Context(), ctx).AbandonTask(id, claims.UserID, req.Reason); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func CompleteTaskHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
		var req types.CompleteTaskReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		if err := task.NewCompleteTaskLogic(r.Context(), ctx).CompleteTask(id, claims.UserID, req.Proof); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}

func ExtendDeadlineHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
		var req types.ExtendDeadlineReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err))
			return
		}
		if err := task.NewExtendDeadlineLogic(r.Context(), ctx).ExtendDeadline(id, claims.UserID, req.NewDeadline); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, nil)
	}
}
