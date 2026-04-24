package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

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
	// todo: add your logic here and delete this line

	return &v1.RejectMedicalRecordResponse{}, nil
}
