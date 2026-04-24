package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyAddRescueRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyAddRescueRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyAddRescueRecordLogic {
	return &ApplyAddRescueRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 申请新增救助记录
func (l *ApplyAddRescueRecordLogic) ApplyAddRescueRecord(in *v1.ApplyAddRescueRecordRequest) (*v1.ApplyAddRescueRecordResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.ApplyAddRescueRecordResponse{}, nil
}
