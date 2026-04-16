// Package logic contains the user service business logic.
// Each exported type corresponds to one RPC method.
package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	userpb "github.com/luyb177/meow-nook/service/user/pb/user"
	"golang.org/x/crypto/bcrypt"
)

// ──────────────────────────────────────────────
// Register
// ──────────────────────────────────────────────

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *RegisterLogic) Register(req *userpb.RegisterReq) (*userpb.RegisterResp, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errorx.ErrBadRequest
	}
	// Hash password before storing.
	_, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeInternalError, "密码加密失败", err)
	}
	// TODO: persist to DB; check uniqueness; return real user ID
	return &userpb.RegisterResp{UserId: 1, Token: "placeholder"}, nil
}

// ──────────────────────────────────────────────
// Login
// ──────────────────────────────────────────────

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *LoginLogic) Login(req *userpb.LoginReq) (*userpb.LoginResp, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errorx.ErrBadRequest
	}
	// TODO: fetch user from DB; compare password hash
	return &userpb.LoginResp{UserId: 1, Token: "placeholder"}, nil
}

// ──────────────────────────────────────────────
// Logout
// ──────────────────────────────────────────────

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *LogoutLogic) Logout(req *userpb.LogoutReq) error {
	// TODO: add token to blocklist in Redis
	return nil
}

// ──────────────────────────────────────────────
// GetUserInfo
// ──────────────────────────────────────────────

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetUserInfoLogic) GetUserInfo(req *userpb.GetUserInfoReq) (*userpb.GetUserInfoResp, error) {
	// TODO: fetch from DB
	return &userpb.GetUserInfoResp{
		User: &userpb.UserInfo{Id: req.UserId},
	}, nil
}

// ──────────────────────────────────────────────
// UpdateUserInfo
// ──────────────────────────────────────────────

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *userpb.UpdateUserInfoReq) error {
	// TODO: update in DB
	return nil
}

// ──────────────────────────────────────────────
// ChangePassword
// ──────────────────────────────────────────────

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ChangePasswordLogic) ChangePassword(req *userpb.ChangePasswordReq) error {
	if req.OldPassword == "" || req.NewPassword == "" {
		return errorx.ErrBadRequest
	}
	// TODO: verify old password and update hash in DB
	return nil
}

// ──────────────────────────────────────────────
// GetPoints
// ──────────────────────────────────────────────

type GetPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPointsLogic {
	return &GetPointsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetPointsLogic) GetPoints(req *userpb.GetPointsReq) (*userpb.GetPointsResp, error) {
	// TODO: fetch from DB
	return &userpb.GetPointsResp{Points: 0}, nil
}

// ──────────────────────────────────────────────
// AddPoints
// ──────────────────────────────────────────────

type AddPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPointsLogic {
	return &AddPointsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AddPointsLogic) AddPoints(req *userpb.AddPointsReq) (*userpb.AddPointsResp, error) {
	// TODO: update points in DB; record log
	return &userpb.AddPointsResp{NewTotal: req.Delta}, nil
}

// ──────────────────────────────────────────────
// ListPointLogs
// ──────────────────────────────────────────────

type ListPointLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPointLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPointLogsLogic {
	return &ListPointLogsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListPointLogsLogic) ListPointLogs(req *userpb.ListPointLogsReq) (*userpb.ListPointLogsResp, error) {
	// TODO: query from DB with pagination
	return &userpb.ListPointLogsResp{Logs: []*userpb.PointLog{}, Total: 0}, nil
}

// ──────────────────────────────────────────────
// GetNotificationSettings
// ──────────────────────────────────────────────

type GetNotificationSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNotificationSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotificationSettingsLogic {
	return &GetNotificationSettingsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetNotificationSettingsLogic) GetNotificationSettings(req *userpb.GetNotificationSettingsReq) (*userpb.GetNotificationSettingsResp, error) {
	// TODO: fetch from DB
	return &userpb.GetNotificationSettingsResp{
		Settings: &userpb.NotificationSettings{
			TaskNotify: true, AdoptionNotify: true, PointsNotify: true,
			SystemNotify: true, MessageNotify: true,
		},
	}, nil
}

// ──────────────────────────────────────────────
// UpdateNotificationSettings
// ──────────────────────────────────────────────

type UpdateNotificationSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateNotificationSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNotificationSettingsLogic {
	return &UpdateNotificationSettingsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateNotificationSettingsLogic) UpdateNotificationSettings(req *userpb.UpdateNotificationSettingsReq) error {
	// TODO: update in DB
	return nil
}

// ──────────────────────────────────────────────
// SubmitFeedback
// ──────────────────────────────────────────────

type SubmitFeedbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitFeedbackLogic {
	return &SubmitFeedbackLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *SubmitFeedbackLogic) SubmitFeedback(req *userpb.SubmitFeedbackReq) (*userpb.SubmitFeedbackResp, error) {
	if req.FeedbackType == "" || req.Content == "" {
		return nil, errorx.ErrBadRequest
	}
	// TODO: persist to DB
	return &userpb.SubmitFeedbackResp{FeedbackId: 1}, nil
}
