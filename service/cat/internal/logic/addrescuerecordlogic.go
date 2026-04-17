package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRescueRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddRescueRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRescueRecordLogic {
	return &AddRescueRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Rescue records
func (l *AddRescueRecordLogic) AddRescueRecord(in *v1.AddRescueRecordReq) (*v1.AddRescueRecordResp, error) {
	// todo: add your logic here and delete this line

	return &v1.AddRescueRecordResp{}, nil
}
