package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCatTaskStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCatTaskStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCatTaskStatusLogic {
	return &UpdateCatTaskStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理员更新任务状态
func (l *UpdateCatTaskStatusLogic) UpdateCatTaskStatus(in *v1.UpdateCatTaskStatusRequest) (*v1.UpdateCatTaskStatusResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.UpdateCatTaskStatusResponse{}, nil
}
