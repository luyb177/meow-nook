package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
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
	// TODO(permission): check that in.OperatorId has admin role

	if _, err := l.svcCtx.Repo.Cat.GetCatTaskByID(l.ctx, in.TaskId); err != nil {
		return nil, errorx.ErrTaskNotFound
	}

	if err := l.svcCtx.Repo.Cat.UpdateCatTask(l.ctx, in.TaskId, map[string]interface{}{
		"status": in.Status,
	}); err != nil {
		logger.Error("UpdateCatTaskStatus: update failed")
		return nil, errorx.WrapInternal("更新任务状态失败", err)
	}

	return &v1.UpdateCatTaskStatusResponse{Success: true}, nil
}
