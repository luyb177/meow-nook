package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyCreateCatTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyCreateCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyCreateCatTaskLogic {
	return &ApplyCreateCatTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 申请创建任务
func (l *ApplyCreateCatTaskLogic) ApplyCreateCatTask(in *v1.ApplyCreateCatTaskRequest) (*v1.ApplyCreateCatTaskResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.ApplyCreateCatTaskResponse{}, nil
}
