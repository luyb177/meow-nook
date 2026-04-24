package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &v1.ApproveMedicalRecordResponse{}, nil
}
