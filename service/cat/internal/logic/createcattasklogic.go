package logic

import (
"context"

"github.com/luyb177/meow-nook/common/errorx"
"github.com/luyb177/meow-nook/common/logger"
"github.com/luyb177/meow-nook/service/cat/internal/repo/task"
"github.com/luyb177/meow-nook/service/cat/internal/svc"
v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
"go.uber.org/zap"

"github.com/zeromicro/go-zero/core/logx"
)

type CreateCatTaskLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewCreateCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCatTaskLogic {
return &CreateCatTaskLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

// 管理员直接创建正式任务（无需审核）
func (l *CreateCatTaskLogic) CreateCatTask(in *v1.CreateCatTaskRequest) (*v1.CreateCatTaskResponse, error) {
if in.CatId == 0 {
return nil, errorx.Wrap(errorx.CodeBadRequest, "cat_id不能为空", errorx.ErrBadRequest)
}
if in.Title == "" {
return nil, errorx.Wrap(errorx.CodeBadRequest, "任务标题不能为空", errorx.ErrBadRequest)
}
if in.OperatorId == 0 {
return nil, errorx.Wrap(errorx.CodeUnauthorized, "操作人ID不能为空", errorx.ErrUnauthorized)
}

urgency := in.UrgencyLevel
if urgency == "" {
urgency = task.UrgencyNormal
}

newTask := &task.CatTask{
CatID:           in.CatId,
CreatorID:       in.OperatorId,
Title:           in.Title,
TaskType:        in.TaskType,
Summary:         in.Summary,
Detail:          in.Detail,
UrgencyLevel:    urgency,
DifficultyLevel: in.DifficultyLevel,
RewardPoints:    in.RewardPoints,
MaxClaimers:     1,
Status:          task.TaskStatusPending,
}

if in.Deadline != nil {
t := in.Deadline.AsTime()
newTask.DeadlineAt = &t
}

if err := l.svcCtx.Repo.Task.CreateTask(l.ctx, newTask); err != nil {
logger.Error("管理员创建任务失败", zap.Error(err))
return nil, errorx.WrapDBInsert("创建任务失败", err)
}

logger.Info("管理员创建任务成功", zap.Uint64("task_id", newTask.ID))

return &v1.CreateCatTaskResponse{TaskId: newTask.ID}, nil
}
