package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	taskrepo "github.com/luyb177/meow-nook/service/cat/internal/repo/task"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
)

type ClaimTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClaimTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimTaskLogic {
	return &ClaimTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClaimTaskLogic) ClaimTask(in *v1.ClaimTaskRequest) (*v1.ClaimTaskResponse, error) {
	logger.Info("ClaimTask called")

	// 先查任务状态
	task, err := l.svcCtx.Repo.Task.GetTaskByID(l.ctx, in.TaskId)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询任务失败", err)
	}
	if task.Status != taskrepo.TaskStatusPending {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "任务不是待认领状态", errorx.ErrBadRequest)
	}

	// 防止重复认领
	exist, _ := l.svcCtx.Repo.Task.GetTaskClaimByTaskAndUser(l.ctx, in.TaskId, in.UserId)
	if exist != nil {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "您已认领该任务", errorx.ErrBadRequest)
	}

	claim := &taskrepo.TaskClaim{
		TaskID: in.TaskId,
		UserID: in.UserId,
	}

	if err := l.svcCtx.Repo.Task.CreateTaskClaim(l.ctx, claim); err != nil {
		return nil, errorx.WrapDBInsert("认领任务失败", err)
	}

	return &v1.ClaimTaskResponse{
		ClaimId: claim.ID,
		Message: "认领成功",
	}, nil
}
