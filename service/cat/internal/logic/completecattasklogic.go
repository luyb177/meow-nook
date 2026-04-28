package logic

import (
"context"
"encoding/json"

"github.com/luyb177/meow-nook/common/errorx"
"github.com/luyb177/meow-nook/common/logger"
"github.com/luyb177/meow-nook/service/cat/internal/svc"
v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
"go.uber.org/zap"

"github.com/zeromicro/go-zero/core/logx"
)

type CompleteCatTaskLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewCompleteCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteCatTaskLogic {
return &CompleteCatTaskLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

// 完成任务
func (l *CompleteCatTaskLogic) CompleteCatTask(in *v1.CompleteCatTaskRequest) (*v1.CompleteCatTaskResponse, error) {
if in.TaskId == 0 {
return nil, errorx.Wrap(errorx.CodeBadRequest, "task_id不能为空", errorx.ErrBadRequest)
}
if in.UserId == 0 {
return nil, errorx.Wrap(errorx.CodeUnauthorized, "用户未登录", errorx.ErrUnauthorized)
}

imageURLsJSON := ""
if len(in.ImageUrls) > 0 {
b, err := json.Marshal(in.ImageUrls)
if err != nil {
return nil, errorx.Wrap(errorx.CodeBadRequest, "图片URL格式错误", err)
}
imageURLsJSON = string(b)
}

if err := l.svcCtx.Repo.Task.CompleteTask(l.ctx, in.TaskId, in.UserId, in.Content, imageURLsJSON); err != nil {
logger.Error("完成任务失败", zap.Uint64("task_id", in.TaskId), zap.Error(err))
return nil, errorx.WrapDBUpdate("完成任务失败", err)
}

logger.Info("完成任务成功", zap.Uint64("task_id", in.TaskId), zap.Uint64("user_id", in.UserId))

return &v1.CompleteCatTaskResponse{Success: true}, nil
}
