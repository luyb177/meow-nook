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
if in.TaskId == 0 {
return nil, errorx.Wrap(errorx.CodeBadRequest, "task_id不能为空", errorx.ErrBadRequest)
}
if in.UserId == 0 {
return nil, errorx.Wrap(errorx.CodeUnauthorized, "用户未登录", errorx.ErrUnauthorized)
}

if err := l.svcCtx.Repo.Task.AbandonTask(l.ctx, in.TaskId, in.UserId, in.Reason); err != nil {
logger.Error("放弃任务失败", zap.Uint64("task_id", in.TaskId), zap.Error(err))
return nil, errorx.WrapDBUpdate("放弃任务失败", err)
}

logger.Info("放弃任务成功", zap.Uint64("task_id", in.TaskId), zap.Uint64("user_id", in.UserId))

return &v1.AbandonCatTaskResponse{Success: true}, nil
}
