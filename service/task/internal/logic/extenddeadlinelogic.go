package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/task/internal/svc"
	"github.com/luyb177/meow-nook/service/task/pb/task/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExtendDeadlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExtendDeadlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExtendDeadlineLogic {
	return &ExtendDeadlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ExtendDeadlineLogic) ExtendDeadline(in *v1.ExtendDeadlineReq) (*v1.ExtendDeadlineResp, error) {
	// todo: add your logic here and delete this line

	return &v1.ExtendDeadlineResp{}, nil
}
