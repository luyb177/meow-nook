package logic

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveMedicalRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApproveMedicalRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveMedicalRecordLogic {
	return &ApproveMedicalRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApproveMedicalRecordLogic) ApproveMedicalRecord(in *v1.ApproveMedicalRecordRequest) (*v1.ApproveMedicalRecordResponse, error) {
	// TODO(permission): check that in.ReviewerId has admin role

	apply, err := l.svcCtx.Repo.Cat.GetCatMedicalApplyByID(l.ctx, in.RecordApplyId)
	if err != nil {
		return nil, errorx.ErrCatNotFound
	}
	if apply.Status != "pending" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "申请状态不是待审核", nil)
	}

	record := &cat.CatMedicalRecord{
		CatID:        apply.CatID,
		MedicalDate:  time.Now(),
		MedicalType:  apply.MedicalType,
		Content:      apply.Content,
		OperatorName: "", // TODO(user): fetch from user service
	}
	if err := l.svcCtx.Repo.Cat.CreateCatMedicalRecord(l.ctx, record); err != nil {
		logger.Error("ApproveMedicalRecord: create record failed")
		return nil, errorx.WrapInternal("创建医疗记录失败", err)
	}

	if err := l.svcCtx.Repo.Cat.UpdateCatMedicalApplyStatus(l.ctx, apply.ID, "approved", "", in.ReviewerId); err != nil {
		logger.Error("ApproveMedicalRecord: update apply status failed")
		return nil, errorx.WrapInternal("更新申请状态失败", err)
	}

	return &v1.ApproveMedicalRecordResponse{Success: true}, nil
}

// ──────────────────────────────────────────────────────────────────────────────

type RejectMedicalRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRejectMedicalRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectMedicalRecordLogic {
	return &RejectMedicalRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RejectMedicalRecordLogic) RejectMedicalRecord(in *v1.RejectMedicalRecordRequest) (*v1.RejectMedicalRecordResponse, error) {
	// TODO(permission): check that in.ReviewerId has admin role

	apply, err := l.svcCtx.Repo.Cat.GetCatMedicalApplyByID(l.ctx, in.RecordApplyId)
	if err != nil {
		return nil, errorx.ErrCatNotFound
	}
	if apply.Status != "pending" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "申请状态不是待审核", nil)
	}

	if err := l.svcCtx.Repo.Cat.UpdateCatMedicalApplyStatus(l.ctx, apply.ID, "rejected", in.Reason, in.ReviewerId); err != nil {
		logger.Error("RejectMedicalRecord: update apply status failed")
		return nil, errorx.WrapInternal("更新申请状态失败", err)
	}

	return &v1.RejectMedicalRecordResponse{Success: true}, nil
}
