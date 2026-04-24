package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type AbandonCatTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAbandonCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AbandonCatTaskLogic {
	return &AbandonCatTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 申请放弃任务
func (l *AbandonCatTaskLogic) AbandonCatTask(in *v1.AbandonCatTaskRequest) (*v1.AbandonCatTaskResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.AbandonCatTaskResponse{}, nil
}
