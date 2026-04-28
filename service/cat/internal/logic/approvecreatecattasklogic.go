package logic

import (
"context"

"github.com/luyb177/meow-nook/common/errorx"
"github.com/luyb177/meow-nook/common/logger"
"github.com/luyb177/meow-nook/service/cat/internal/repo/task"
"github.com/luyb177/meow-nook/service/cat/internal/svc"
v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
"go.uber.org/zap"
"gorm.io/gorm"

"github.com/zeromicro/go-zero/core/logx"
)

type ApproveCreateCatTaskLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewApproveCreateCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveCreateCatTaskLogic {
return &ApproveCreateCatTaskLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

// 管理员审核志愿者任务申请
func (l *ApproveCreateCatTaskLogic) ApproveCreateCatTask(in *v1.ApproveCreateCatTaskRequest) (*v1.ApproveCreateCatTaskResponse, error) {
if in.ApplyId == 0 {
return nil, errorx.Wrap(errorx.CodeBadRequest, "apply_id不能为空", errorx.ErrBadRequest)
}
if in.OperatorId == 0 {
return nil, errorx.Wrap(errorx.CodeUnauthorized, "操作人ID不能为空", errorx.ErrUnauthorized)
}

var taskID uint64

err := l.svcCtx.Repo.WithTx(l.ctx, func(tx *gorm.DB) error {
// 1. 查询并锁定申请
apply, err := l.svcCtx.Repo.Task.GetTaskApplyByIDForUpdate(l.ctx, in.ApplyId, tx)
if err != nil {
return errorx.WrapDBQuery("查询任务申请失败", err)
}

if apply.Status != task.TaskApplyStatusPending {
return errorx.Wrap(errorx.CodeBadRequest, "申请状态不是待审核", errorx.ErrBadRequest)
}

// 2. 创建正式任务
newTask := &task.CatTask{
CatID:           apply.CatID,
ApplyID:         apply.ID,
CreatorID:       in.OperatorId,
Title:           apply.Title,
TaskType:        apply.TaskType,
Summary:         apply.Summary,
Detail:          apply.Detail,
UrgencyLevel:    in.UrgencyLevel,
DifficultyLevel: in.DifficultyLevel,
RewardPoints:    in.RewardPoints,
MaxClaimers:     1,
Status:          task.TaskStatusPending,
}
if newTask.UrgencyLevel == "" {
newTask.UrgencyLevel = task.UrgencyNormal
}
if apply.DeadlineAt != nil {
newTask.DeadlineAt = apply.DeadlineAt
}

if err := l.svcCtx.Repo.Task.CreateTask(l.ctx, newTask, tx); err != nil {
return errorx.WrapDBInsert("创建正式任务失败", err)
}
taskID = newTask.ID

// 3. 更新申请状态
if err := l.svcCtx.Repo.Task.ApproveTaskApply(
l.ctx,
in.ApplyId,
in.OperatorId,
in.UrgencyLevel,
in.DifficultyLevel,
in.RewardPoints,
taskID,
tx,
); err != nil {
return errorx.WrapDBUpdate("更新申请状态失败", err)
}

return nil
})

if err != nil {
logger.Error("审核任务申请失败", zap.Uint64("apply_id", in.ApplyId), zap.Error(err))
return nil, err
}

logger.Info("审核任务申请成功", zap.Uint64("apply_id", in.ApplyId), zap.Uint64("task_id", taskID))

return &v1.ApproveCreateCatTaskResponse{TaskId: taskID}, nil
}
