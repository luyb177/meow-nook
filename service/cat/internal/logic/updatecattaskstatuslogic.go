package logic

import (
"context"

"github.com/luyb177/meow-nook/common/errorx"
"github.com/luyb177/meow-nook/common/logger"
"github.com/luyb177/meow-nook/service/cat/internal/svc"
v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
"go.uber.org/zap"

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

func (l *UpdateCatTaskStatusLogic) UpdateCatTaskStatus(in *v1.UpdateCatTaskStatusRequest) (*v1.UpdateCatTaskStatusResponse, error) {
if in.TaskId == 0 {
return nil, errorx.Wrap(errorx.CodeBadRequest, "task_id不能为空", errorx.ErrBadRequest)
}
if in.Status == "" {
return nil, errorx.Wrap(errorx.CodeBadRequest, "status不能为空", errorx.ErrBadRequest)
}

if err := l.svcCtx.Repo.Task.UpdateTaskStatus(l.ctx, in.TaskId, in.Status, in.Remark); err != nil {
logger.Error("更新任务状态失败", zap.Uint64("task_id", in.TaskId), zap.Error(err))
return nil, errorx.WrapDBUpdate("更新任务状态失败", err)
}

logger.Info("更新任务状态成功", zap.Uint64("task_id", in.TaskId), zap.String("status", in.Status))

return &v1.UpdateCatTaskStatusResponse{Success: true}, nil
}
