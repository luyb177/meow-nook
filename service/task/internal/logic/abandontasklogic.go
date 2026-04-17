package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/task/internal/svc"
	"github.com/luyb177/meow-nook/service/task/pb/task/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type AbandonTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAbandonTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AbandonTaskLogic {
	return &AbandonTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AbandonTaskLogic) AbandonTask(in *v1.AbandonTaskReq) (*v1.AbandonTaskResp, error) {
	// todo: add your logic here and delete this line

	return &v1.AbandonTaskResp{}, nil
}
