package handler

import (
	"net/http"

	gatewayMiddleware "github.com/luyb177/meow-nook/service/gateway/internal/middleware"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	adoptionHandler "github.com/luyb177/meow-nook/service/gateway/internal/handler/adoption"
	catHandler "github.com/luyb177/meow-nook/service/gateway/internal/handler/cat"
	postHandler "github.com/luyb177/meow-nook/service/gateway/internal/handler/post"
	taskHandler "github.com/luyb177/meow-nook/service/gateway/internal/handler/task"
	userHandler "github.com/luyb177/meow-nook/service/gateway/internal/handler/user"
	"github.com/zeromicro/go-zero/rest"
)

// RegisterHandlers registers all routes on the given server.
func RegisterHandlers(server *rest.Server, ctx *svc.ServiceContext) {
	authMW := gatewayMiddleware.NewAuthMiddleware(ctx.Config.Auth.AccessSecret)

	// ── Auth (no JWT) ──────────────────────────────────────────────────────
	server.AddRoutes([]rest.Route{
		{Method: http.MethodPost, Path: "/api/v1/auth/register", Handler: userHandler.RegisterHandler(ctx)},
		{Method: http.MethodPost, Path: "/api/v1/auth/login", Handler: userHandler.LoginHandler(ctx)},
	})

	// ── User (JWT required) ────────────────────────────────────────────────
	server.AddRoutes(rest.WithMiddlewares(
		[]rest.Middleware{authMW.Handle},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/user/info", Handler: userHandler.GetUserInfoHandler(ctx)},
		rest.Route{Method: http.MethodPut, Path: "/api/v1/user/info", Handler: userHandler.UpdateUserInfoHandler(ctx)},
		rest.Route{Method: http.MethodPut, Path: "/api/v1/user/password", Handler: userHandler.ChangePasswordHandler(ctx)},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/user/points", Handler: userHandler.GetPointsHandler(ctx)},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/user/points/logs", Handler: userHandler.ListPointLogsHandler(ctx)},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/user/notifications", Handler: userHandler.GetNotificationSettingsHandler(ctx)},
		rest.Route{Method: http.MethodPut, Path: "/api/v1/user/notifications", Handler: userHandler.UpdateNotificationSettingsHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/user/feedback", Handler: userHandler.SubmitFeedbackHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/user/logout", Handler: userHandler.LogoutHandler(ctx)},
	))

	// ── Cats ───────────────────────────────────────────────────────────────
	server.AddRoutes(rest.WithMiddlewares(
		[]rest.Middleware{authMW.Handle},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/cats", Handler: catHandler.ListCatsHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/cats", Handler: catHandler.CreateCatHandler(ctx)},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/cats/stats", Handler: catHandler.GetCatStatsHandler(ctx)},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/cats/heatmap", Handler: catHandler.GetHeatmapHandler(ctx)},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/cats/:id", Handler: catHandler.GetCatHandler(ctx)},
		rest.Route{Method: http.MethodPut, Path: "/api/v1/cats/:id", Handler: catHandler.UpdateCatHandler(ctx)},
		rest.Route{Method: http.MethodDelete, Path: "/api/v1/cats/:id", Handler: catHandler.DeleteCatHandler(ctx)},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/cats/:id/rescue", Handler: catHandler.ListRescueRecordsHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/cats/:id/rescue", Handler: catHandler.AddRescueRecordHandler(ctx)},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/cats/:id/health", Handler: catHandler.ListHealthRecordsHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/cats/:id/health", Handler: catHandler.AddHealthRecordHandler(ctx)},
	))

	// ── Tasks ──────────────────────────────────────────────────────────────
	server.AddRoutes(rest.WithMiddlewares(
		[]rest.Middleware{authMW.Handle},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/tasks", Handler: taskHandler.ListTasksHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/tasks", Handler: taskHandler.CreateTaskHandler(ctx)},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/tasks/:id", Handler: taskHandler.GetTaskHandler(ctx)},
		rest.Route{Method: http.MethodPut, Path: "/api/v1/tasks/:id", Handler: taskHandler.UpdateTaskHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/tasks/:id/claim", Handler: taskHandler.ClaimTaskHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/tasks/:id/abandon", Handler: taskHandler.AbandonTaskHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/tasks/:id/complete", Handler: taskHandler.CompleteTaskHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/tasks/:id/extend", Handler: taskHandler.ExtendDeadlineHandler(ctx)},
	))

	// ── Adoption ───────────────────────────────────────────────────────────
	server.AddRoutes(rest.WithMiddlewares(
		[]rest.Middleware{authMW.Handle},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/adoption", Handler: adoptionHandler.ListAdoptionsHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/adoption/apply", Handler: adoptionHandler.ApplyAdoptionHandler(ctx)},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/adoption/:id", Handler: adoptionHandler.GetAdoptionHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/adoption/:id/review", Handler: adoptionHandler.ReviewAdoptionHandler(ctx)},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/adoption/:id/followup", Handler: adoptionHandler.ListFollowUpsHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/adoption/:id/followup", Handler: adoptionHandler.SubmitFollowUpHandler(ctx)},
	))

	// ── Posts ──────────────────────────────────────────────────────────────
	server.AddRoutes(rest.WithMiddlewares(
		[]rest.Middleware{authMW.Handle},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/posts", Handler: postHandler.ListPostsHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/posts", Handler: postHandler.CreatePostHandler(ctx)},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/posts/:id", Handler: postHandler.GetPostHandler(ctx)},
		rest.Route{Method: http.MethodDelete, Path: "/api/v1/posts/:id", Handler: postHandler.DeletePostHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/posts/:id/like", Handler: postHandler.LikePostHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/posts/:id/unlike", Handler: postHandler.UnlikePostHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/posts/:id/favorite", Handler: postHandler.FavoritePostHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/posts/:id/unfavorite", Handler: postHandler.UnfavoritePostHandler(ctx)},
		rest.Route{Method: http.MethodGet, Path: "/api/v1/posts/:id/comments", Handler: postHandler.ListCommentsHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/posts/:id/comments", Handler: postHandler.AddCommentHandler(ctx)},
		rest.Route{Method: http.MethodDelete, Path: "/api/v1/posts/:id/comments/:cid", Handler: postHandler.DeleteCommentHandler(ctx)},
		rest.Route{Method: http.MethodPost, Path: "/api/v1/posts/:id/report", Handler: postHandler.ReportPostHandler(ctx)},
	))
}
