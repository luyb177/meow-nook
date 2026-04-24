package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &v1.ApplyAddMedicalRecordResponse{}, nil
}
