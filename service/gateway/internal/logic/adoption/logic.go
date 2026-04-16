package adoption

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type ListAdoptionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAdoptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAdoptionsLogic {
	return &ListAdoptionsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListAdoptionsLogic) ListAdoptions(userID int64, req *types.ListAdoptionReq) (*types.ListAdoptionResp, error) {
	// TODO: call adoption gRPC service
	return &types.ListAdoptionResp{Applications: []types.AdoptionApplication{}, Total: 0}, nil
}

// ──────────────────────────────────────────────

type ApplyAdoptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyAdoptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyAdoptionLogic {
	return &ApplyAdoptionLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ApplyAdoptionLogic) ApplyAdoption(userID int64, req *types.ApplyAdoptionReq) (*types.ApplyAdoptionResp, error) {
	if req.CatId == 0 {
		return nil, errorx.New(errorx.CodeBadRequest, "必须指定领养的猫咪")
	}
	// TODO: call adoption gRPC service
	return &types.ApplyAdoptionResp{ApplicationId: 1}, nil
}

// ──────────────────────────────────────────────

type GetAdoptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdoptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdoptionLogic {
	return &GetAdoptionLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAdoptionLogic) GetAdoption(applicationID int64) (*types.AdoptionApplication, error) {
	// TODO: call adoption gRPC service
	return &types.AdoptionApplication{Id: applicationID}, nil
}

// ──────────────────────────────────────────────

type ReviewAdoptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReviewAdoptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReviewAdoptionLogic {
	return &ReviewAdoptionLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ReviewAdoptionLogic) ReviewAdoption(applicationID, reviewerID int64, req *types.ReviewApplicationReq) error {
	if !req.Approved && req.RejectReason == "" {
		return errorx.New(errorx.CodeBadRequest, "拒绝时必须填写原因")
	}
	// TODO: call adoption gRPC service
	return nil
}

// ──────────────────────────────────────────────

type ListFollowUpsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListFollowUpsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFollowUpsLogic {
	return &ListFollowUpsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListFollowUpsLogic) ListFollowUps(applicationID int64) ([]types.FollowUpRecord, error) {
	// TODO: call adoption gRPC service
	return []types.FollowUpRecord{}, nil
}

// ──────────────────────────────────────────────

type SubmitFollowUpLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitFollowUpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitFollowUpLogic {
	return &SubmitFollowUpLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *SubmitFollowUpLogic) SubmitFollowUp(applicationID, userID int64, req *types.SubmitFollowUpReq) error {
	if req.Stage == "" || req.Content == "" {
		return errorx.New(errorx.CodeBadRequest, "回访阶段和内容不能为空")
	}
	// TODO: call adoption gRPC service
	return nil
}
