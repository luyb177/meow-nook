// Package server provides the gRPC server implementation for the user service.
// The actual business logic is in the logic sub-package; this layer only
// translates between the gRPC wire format and the internal types.
package server

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/user/internal/logic"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	userpb "github.com/luyb177/meow-nook/service/user/pb/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterUserServer wires up the UserServiceServer to the given gRPC server.
func RegisterUserServer(s *grpc.Server, ctx *svc.ServiceContext) {
	userpb.RegisterUserServiceServer(s, &UserServer{ctx: ctx})
}

// UserServer implements userpb.UserServiceServer.
type UserServer struct {
	userpb.UnimplementedUserServiceServer
	ctx *svc.ServiceContext
}

func toGRPCError(err error) error {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*errorx.AppError); ok {
		switch ae.Code {
		case errorx.CodeNotFound, errorx.CodeUserNotFound, errorx.CodeCatNotFound,
			errorx.CodeTaskNotFound, errorx.CodeAdoptionNotFound, errorx.CodePostNotFound:
			return status.Error(codes.NotFound, ae.Msg)
		case errorx.CodeUnauthorized, errorx.CodeTokenInvalid, errorx.CodeTokenExpired:
			return status.Error(codes.Unauthenticated, ae.Msg)
		case errorx.CodeForbidden, errorx.CodePermissionDenied, errorx.CodeInsufficientPoints:
			return status.Error(codes.PermissionDenied, ae.Msg)
		case errorx.CodeConflict, errorx.CodeUserAlreadyExists,
			errorx.CodeTaskAlreadyClaimed, errorx.CodeTaskFull, errorx.CodeAdoptionAlreadyApplied:
			return status.Error(codes.AlreadyExists, ae.Msg)
		case errorx.CodeBadRequest, errorx.CodeUnprocessable:
			return status.Error(codes.InvalidArgument, ae.Msg)
		default:
			return status.Error(codes.Internal, ae.Msg)
		}
	}
	return status.Error(codes.Internal, "服务器内部错误")
}

// ──────────────────────────────────────────────
// Auth
// ──────────────────────────────────────────────

func (s *UserServer) Register(ctx context.Context, req *userpb.RegisterReq) (*userpb.RegisterResp, error) {
	l := logic.NewRegisterLogic(ctx, s.ctx)
	resp, err := l.Register(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *UserServer) Login(ctx context.Context, req *userpb.LoginReq) (*userpb.LoginResp, error) {
	l := logic.NewLoginLogic(ctx, s.ctx)
	resp, err := l.Login(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *UserServer) Logout(ctx context.Context, req *userpb.LogoutReq) (*userpb.LogoutResp, error) {
	l := logic.NewLogoutLogic(ctx, s.ctx)
	if err := l.Logout(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &userpb.LogoutResp{}, nil
}

// ──────────────────────────────────────────────
// Profile
// ──────────────────────────────────────────────

func (s *UserServer) GetUserInfo(ctx context.Context, req *userpb.GetUserInfoReq) (*userpb.GetUserInfoResp, error) {
	l := logic.NewGetUserInfoLogic(ctx, s.ctx)
	resp, err := l.GetUserInfo(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *UserServer) UpdateUserInfo(ctx context.Context, req *userpb.UpdateUserInfoReq) (*userpb.UpdateUserInfoResp, error) {
	l := logic.NewUpdateUserInfoLogic(ctx, s.ctx)
	if err := l.UpdateUserInfo(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &userpb.UpdateUserInfoResp{}, nil
}

func (s *UserServer) ChangePassword(ctx context.Context, req *userpb.ChangePasswordReq) (*userpb.ChangePasswordResp, error) {
	l := logic.NewChangePasswordLogic(ctx, s.ctx)
	if err := l.ChangePassword(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &userpb.ChangePasswordResp{}, nil
}

// ──────────────────────────────────────────────
// Points
// ──────────────────────────────────────────────

func (s *UserServer) GetPoints(ctx context.Context, req *userpb.GetPointsReq) (*userpb.GetPointsResp, error) {
	l := logic.NewGetPointsLogic(ctx, s.ctx)
	resp, err := l.GetPoints(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *UserServer) AddPoints(ctx context.Context, req *userpb.AddPointsReq) (*userpb.AddPointsResp, error) {
	l := logic.NewAddPointsLogic(ctx, s.ctx)
	resp, err := l.AddPoints(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *UserServer) ListPointLogs(ctx context.Context, req *userpb.ListPointLogsReq) (*userpb.ListPointLogsResp, error) {
	l := logic.NewListPointLogsLogic(ctx, s.ctx)
	resp, err := l.ListPointLogs(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ──────────────────────────────────────────────
// Notifications
// ──────────────────────────────────────────────

func (s *UserServer) GetNotificationSettings(ctx context.Context, req *userpb.GetNotificationSettingsReq) (*userpb.GetNotificationSettingsResp, error) {
	l := logic.NewGetNotificationSettingsLogic(ctx, s.ctx)
	resp, err := l.GetNotificationSettings(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *UserServer) UpdateNotificationSettings(ctx context.Context, req *userpb.UpdateNotificationSettingsReq) (*userpb.UpdateNotificationSettingsResp, error) {
	l := logic.NewUpdateNotificationSettingsLogic(ctx, s.ctx)
	if err := l.UpdateNotificationSettings(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &userpb.UpdateNotificationSettingsResp{}, nil
}

// ──────────────────────────────────────────────
// Feedback
// ──────────────────────────────────────────────

func (s *UserServer) SubmitFeedback(ctx context.Context, req *userpb.SubmitFeedbackReq) (*userpb.SubmitFeedbackResp, error) {
	l := logic.NewSubmitFeedbackLogic(ctx, s.ctx)
	resp, err := l.SubmitFeedback(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}
