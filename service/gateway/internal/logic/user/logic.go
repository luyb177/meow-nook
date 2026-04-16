package user

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luyb177/meow-nook/common/errorx"
	commonMiddleware "github.com/luyb177/meow-nook/common/middleware"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (*types.RegisterResp, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errorx.ErrBadRequest
	}
	// TODO: call user gRPC service to create user
	// Placeholder implementation
	userID := int64(1)
	token, err := l.generateToken(userID, req.Username, "user")
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeInternalError, "生成令牌失败", err)
	}
	return &types.RegisterResp{UserId: userID, Token: token}, nil
}

func (l *RegisterLogic) generateToken(userID int64, username, role string) (string, error) {
	claims := commonMiddleware.Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(
				time.Duration(l.svcCtx.Config.Auth.AccessExpire) * time.Second,
			)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))
}

// ──────────────────────────────────────────────

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errorx.ErrBadRequest
	}
	// TODO: call user gRPC service to verify credentials
	// Placeholder implementation
	userID := int64(1)
	claims := commonMiddleware.Claims{
		UserID:   userID,
		Username: req.Username,
		Role:     "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(
				time.Duration(l.svcCtx.Config.Auth.AccessExpire) * time.Second,
			)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeInternalError, "生成令牌失败", err)
	}
	return &types.LoginResp{UserId: userID, Token: tokenStr}, nil
}

// ──────────────────────────────────────────────

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *LogoutLogic) Logout(userID int64) error {
	// TODO: call user gRPC service to invalidate token (e.g. add to blocklist)
	return nil
}

// ──────────────────────────────────────────────

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetUserInfoLogic) GetUserInfo(userID int64) (*types.UserInfoResp, error) {
	// TODO: call user gRPC service
	return &types.UserInfoResp{Id: userID}, nil
}

// ──────────────────────────────────────────────

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(userID int64, req *types.UpdateUserInfoReq) error {
	// TODO: call user gRPC service
	return nil
}

// ──────────────────────────────────────────────

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ChangePasswordLogic) ChangePassword(userID int64, req *types.ChangePasswordReq) error {
	if req.OldPassword == "" || req.NewPassword == "" {
		return errorx.ErrBadRequest
	}
	// TODO: call user gRPC service
	return nil
}

// ──────────────────────────────────────────────

type GetPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPointsLogic {
	return &GetPointsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetPointsLogic) GetPoints(userID int64) (*types.PointsResp, error) {
	// TODO: call user gRPC service
	return &types.PointsResp{Points: 0}, nil
}

// ──────────────────────────────────────────────

type ListPointLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPointLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPointLogsLogic {
	return &ListPointLogsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListPointLogsLogic) ListPointLogs(userID int64, req *types.PageReq) (*types.PointLogsResp, error) {
	// TODO: call user gRPC service
	return &types.PointLogsResp{Logs: []types.PointLog{}, Total: 0}, nil
}

// ──────────────────────────────────────────────

type GetNotificationSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNotificationSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotificationSettingsLogic {
	return &GetNotificationSettingsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetNotificationSettingsLogic) GetNotificationSettings(userID int64) (*types.NotificationSettings, error) {
	// TODO: call user gRPC service
	return &types.NotificationSettings{
		TaskNotify:     true,
		AdoptionNotify: true,
		PointsNotify:   true,
		SystemNotify:   true,
		MessageNotify:  true,
	}, nil
}

// ──────────────────────────────────────────────

type UpdateNotificationSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateNotificationSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNotificationSettingsLogic {
	return &UpdateNotificationSettingsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateNotificationSettingsLogic) UpdateNotificationSettings(userID int64, req *types.NotificationSettings) error {
	// TODO: call user gRPC service
	return nil
}

// ──────────────────────────────────────────────

type SubmitFeedbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitFeedbackLogic {
	return &SubmitFeedbackLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *SubmitFeedbackLogic) SubmitFeedback(userID int64, req *types.FeedbackReq) (*types.FeedbackResp, error) {
	if req.FeedbackType == "" || req.Content == "" {
		return nil, errorx.ErrBadRequest
	}
	// TODO: call user gRPC service
	return &types.FeedbackResp{FeedbackId: 1}, nil
}
