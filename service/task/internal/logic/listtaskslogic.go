package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/task/internal/svc"
	"github.com/luyb177/meow-nook/service/task/pb/task/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTasksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTasksLogic {
	return &ListTasksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListTasksLogic) ListTasks(in *v1.ListTasksReq) (*v1.ListTasksResp, error) {
	// todo: add your logic here and delete this line

	return &v1.ListTasksResp{}, nil
}
