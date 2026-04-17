package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddHealthRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddHealthRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddHealthRecordLogic {
	return &AddHealthRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Health records
func (l *AddHealthRecordLogic) AddHealthRecord(in *v1.AddHealthRecordReq) (*v1.AddHealthRecordResp, error) {
	// todo: add your logic here and delete this line

	return &v1.AddHealthRecordResp{}, nil
}
