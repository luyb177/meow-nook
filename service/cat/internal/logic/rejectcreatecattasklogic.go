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

type RejectCreateCatTaskLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewRejectCreateCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectCreateCatTaskLogic {
return &RejectCreateCatTaskLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

func (l *RejectCreateCatTaskLogic) RejectCreateCatTask(in *v1.RejectCreateCatTaskRequest) (*v1.RejectCreateCatTaskResponse, error) {
if in.ApplyId == 0 {
return nil, errorx.Wrap(errorx.CodeBadRequest, "apply_id不能为空", errorx.ErrBadRequest)
}
if in.Reason == "" {
return nil, errorx.Wrap(errorx.CodeBadRequest, "驳回原因不能为空", errorx.ErrBadRequest)
}

if err := l.svcCtx.Repo.Task.RejectTaskApply(l.ctx, in.ApplyId, in.OperatorId, in.Reason); err != nil {
logger.Error("驳回任务申请失败", zap.Uint64("apply_id", in.ApplyId), zap.Error(err))
return nil, errorx.WrapDBUpdate("驳回任务申请失败", err)
}

logger.Info("驳回任务申请成功", zap.Uint64("apply_id", in.ApplyId))

return &v1.RejectCreateCatTaskResponse{Success: true}, nil
}
