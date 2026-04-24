package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCatTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCatTaskLogic {
	return &CreateCatTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理员直接创建正式任务（无需审核）
func (l *CreateCatTaskLogic) CreateCatTask(in *v1.CreateCatTaskRequest) (*v1.CreateCatTaskResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.CreateCatTaskResponse{}, nil
}
