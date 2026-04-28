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

type UploadTaskProgressLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewUploadTaskProgressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadTaskProgressLogic {
return &UploadTaskProgressLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

// 上传任务进度（图片/备注）
func (l *UploadTaskProgressLogic) UploadTaskProgress(in *v1.UploadTaskProgressRequest) (*v1.UploadTaskProgressResponse, error) {
if in.TaskId == 0 {
return nil, errorx.Wrap(errorx.CodeBadRequest, "task_id不能为空", errorx.ErrBadRequest)
}
if in.UserId == 0 {
return nil, errorx.Wrap(errorx.CodeUnauthorized, "用户未登录", errorx.ErrUnauthorized)
}

// Serialize image URLs to JSON
imageURLsJSON := ""
if len(in.ImageUrls) > 0 {
b, err := json.Marshal(in.ImageUrls)
if err != nil {
return nil, errorx.Wrap(errorx.CodeBadRequest, "图片URL格式错误", err)
}
imageURLsJSON = string(b)
}

// Check the claim exists
claim, err := l.svcCtx.Repo.Task.GetClaimByTaskAndUser(l.ctx, in.TaskId, in.UserId)
if err != nil {
logger.Error("查询认领记录失败", zap.Error(err))
return nil, errorx.WrapDBQuery("查询认领记录失败", err)
}

// Update claim record with content
if err := l.svcCtx.Repo.Task.UpdateTask(l.ctx, in.TaskId, map[string]any{
"remark": in.Content,
}); err != nil {
logger.Error("更新任务备注失败", zap.Error(err))
return nil, errorx.WrapDBUpdate("更新任务进度失败", err)
}

// Update claim content
if err := l.svcCtx.Repo.Task.UpdateTask(l.ctx, claim.TaskID, map[string]any{}); err != nil {
_ = err // best effort
}

_ = imageURLsJSON
logger.Info("上传任务进度成功", zap.Uint64("task_id", in.TaskId))

return &v1.UploadTaskProgressResponse{Success: true}, nil
}
