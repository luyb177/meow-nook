package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyAddMedicalRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyAddMedicalRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyAddMedicalRecordLogic {
	return &ApplyAddMedicalRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 申请新增医疗记录
func (l *ApplyAddMedicalRecordLogic) ApplyAddMedicalRecord(in *v1.ApplyAddMedicalRecordRequest) (*v1.ApplyAddMedicalRecordResponse, error) {
	if _, err := l.svcCtx.Repo.Cat.GetCatByID(l.ctx, in.CatId); err != nil {
		return nil, errorx.ErrCatNotFound
	}

	apply := &cat.CatMedicalApply{
		CatID:           in.CatId,
		MedicalType:     in.MedicalType,
		Content:         in.Content,
		ApplicantUserID: in.OperatorId,
		Status:          "pending",
	}

	if err := l.svcCtx.Repo.Cat.CreateCatMedicalApply(l.ctx, apply); err != nil {
		logger.Error("ApplyAddMedicalRecord: create apply failed")
		return nil, errorx.WrapInternal("创建医疗记录申请失败", err)
	}

	return &v1.ApplyAddMedicalRecordResponse{RecordApplyId: apply.ID}, nil
}
