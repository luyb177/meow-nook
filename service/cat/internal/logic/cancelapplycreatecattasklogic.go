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

type CancelApplyCreateCatTaskLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewCancelApplyCreateCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelApplyCreateCatTaskLogic {
return &CancelApplyCreateCatTaskLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

// 志愿者取消申请创建任务
func (l *CancelApplyCreateCatTaskLogic) CancelApplyCreateCatTask(in *v1.CancelApplyCreateCatTaskRequest) (*v1.Response, error) {
if in.ApplyId == 0 {
return nil, errorx.Wrap(errorx.CodeBadRequest, "apply_id不能为空", errorx.ErrBadRequest)
}
if in.UserId == 0 {
return nil, errorx.Wrap(errorx.CodeUnauthorized, "用户未登录", errorx.ErrUnauthorized)
}

if err := l.svcCtx.Repo.Task.CancelTaskApply(l.ctx, in.ApplyId, in.UserId, in.Reason); err != nil {
logger.Error("取消任务申请失败", zap.Uint64("apply_id", in.ApplyId), zap.Error(err))
return nil, errorx.WrapDBUpdate("取消任务申请失败", err)
}

logger.Info("取消任务申请成功", zap.Uint64("apply_id", in.ApplyId))

return &v1.Response{}, nil
}
