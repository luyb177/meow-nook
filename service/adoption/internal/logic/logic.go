package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/adoption/internal/svc"
	adoptionpb "github.com/luyb177/meow-nook/service/adoption/pb/adoption"
)

type ApplyAdoptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyAdoptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyAdoptionLogic {
	return &ApplyAdoptionLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ApplyAdoptionLogic) ApplyAdoption(req *adoptionpb.ApplyAdoptionReq) (*adoptionpb.ApplyAdoptionResp, error) {
	if req.CatId == 0 {
		return nil, errorx.New(errorx.CodeBadRequest, "必须指定领养的猫咪")
	}
	// TODO: check cat adoption status; check user points; persist to DB
	return &adoptionpb.ApplyAdoptionResp{ApplicationId: 1}, nil
}

type ReviewApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReviewApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReviewApplicationLogic {
	return &ReviewApplicationLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ReviewApplicationLogic) ReviewApplication(req *adoptionpb.ReviewApplicationReq) error {
	if !req.Approved && req.RejectReason == "" {
		return errorx.New(errorx.CodeBadRequest, "拒绝时必须填写原因")
	}
	// TODO: update status in DB; notify applicant
	return nil
}

type GetApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApplicationLogic {
	return &GetApplicationLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetApplicationLogic) GetApplication(req *adoptionpb.GetApplicationReq) (*adoptionpb.GetApplicationResp, error) {
	// TODO: fetch from DB
	return &adoptionpb.GetApplicationResp{
		Application: &adoptionpb.AdoptionApplication{Id: req.ApplicationId},
	}, nil
}

type ListApplicationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListApplicationsLogic {
	return &ListApplicationsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListApplicationsLogic) ListApplications(req *adoptionpb.ListApplicationsReq) (*adoptionpb.ListApplicationsResp, error) {
	// TODO: query from DB with filters
	return &adoptionpb.ListApplicationsResp{
		Applications: []*adoptionpb.AdoptionApplication{},
		Total:        0,
	}, nil
}

type SubmitFollowUpLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitFollowUpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitFollowUpLogic {
	return &SubmitFollowUpLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *SubmitFollowUpLogic) SubmitFollowUp(req *adoptionpb.SubmitFollowUpReq) (*adoptionpb.SubmitFollowUpResp, error) {
	if req.Stage == "" || req.Content == "" {
		return nil, errorx.New(errorx.CodeBadRequest, "回访阶段和内容不能为空")
	}
	// TODO: persist to DB
	return &adoptionpb.SubmitFollowUpResp{RecordId: 1}, nil
}

type ListFollowUpsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListFollowUpsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFollowUpsLogic {
	return &ListFollowUpsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListFollowUpsLogic) ListFollowUps(req *adoptionpb.ListFollowUpsReq) (*adoptionpb.ListFollowUpsResp, error) {
	// TODO: query from DB
	return &adoptionpb.ListFollowUpsResp{Records: []*adoptionpb.FollowUpRecord{}}, nil
}
