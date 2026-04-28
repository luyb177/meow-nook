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

type ClaimCatTaskLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewClaimCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimCatTaskLogic {
return &ClaimCatTaskLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

// 申请认领任务
func (l *ClaimCatTaskLogic) ClaimCatTask(in *v1.ClaimCatTaskRequest) (*v1.ClaimCatTaskResponse, error) {
if in.TaskId == 0 {
return nil, errorx.Wrap(errorx.CodeBadRequest, "task_id不能为空", errorx.ErrBadRequest)
}
if in.UserId == 0 {
return nil, errorx.Wrap(errorx.CodeUnauthorized, "用户未登录", errorx.ErrUnauthorized)
}

if err := l.svcCtx.Repo.Task.ClaimTask(l.ctx, in.TaskId, in.UserId); err != nil {
logger.Error("认领任务失败", zap.Uint64("task_id", in.TaskId), zap.Error(err))
return nil, errorx.WrapDBUpdate("认领任务失败", err)
}

logger.Info("认领任务成功", zap.Uint64("task_id", in.TaskId), zap.Uint64("user_id", in.UserId))

return &v1.ClaimCatTaskResponse{Success: true}, nil
}
