package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveCreateCatTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApproveCreateCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveCreateCatTaskLogic {
	return &ApproveCreateCatTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理员审核志愿者任务申请
func (l *ApproveCreateCatTaskLogic) ApproveCreateCatTask(in *v1.ApproveCreateCatTaskRequest) (*v1.ApproveCreateCatTaskResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.ApproveCreateCatTaskResponse{}, nil
}
