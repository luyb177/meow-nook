package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/task/internal/svc"
	"github.com/luyb177/meow-nook/service/task/pb/task/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskLogic {
	return &UpdateTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateTaskLogic) UpdateTask(in *v1.UpdateTaskReq) (*v1.UpdateTaskResp, error) {
	// todo: add your logic here and delete this line

	return &v1.UpdateTaskResp{}, nil
}
