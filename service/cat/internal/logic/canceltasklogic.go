package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	taskrepo "github.com/luyb177/meow-nook/service/cat/internal/repo/task"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelTaskLogic {
	return &CancelTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelTaskLogic) CancelTask(in *v1.CancelTaskRequest) (*v1.CancelTaskResponse, error) {
	logger.Info("CancelTask called")

	if in.TaskId == 0 || in.OperatorId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "task_id and operator_id are required", errorx.ErrBadRequest)
	}

	// 可选：校验操作人权限（管理员或创建人）
	task, err := l.svcCtx.Repo.Task.GetTaskByID(l.ctx, in.TaskId)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询任务失败", err)
	}
	if task.Status == taskrepo.TaskStatusCompleted {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "task is completed", errorx.ErrBadRequest)
	}

	if err := l.svcCtx.Repo.Task.CancelTask(l.ctx, in.TaskId, in.OperatorId, in.Reason); err != nil {
		return nil, errorx.WrapDBUpdate("取消任务失败", err)
	}

	// 记录流转日志
	flow := &taskrepo.TaskFlow{
		TaskID:     in.TaskId,
		UserID:     in.OperatorId,
		Action:     "cancel",
		FromStatus: task.Status,
		ToStatus:   taskrepo.TaskStatusCancelled,
		Remark:     in.Reason,
	}
	_ = l.svcCtx.Repo.Task.CreateTaskFlow(l.ctx, flow)

	return &v1.CancelTaskResponse{
		Message: "任务已取消",
	}, nil
}
