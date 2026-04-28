package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/task"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyCreateCatTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyCreateCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyCreateCatTaskLogic {
	return &ApplyCreateCatTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 申请创建任务
func (l *ApplyCreateCatTaskLogic) ApplyCreateCatTask(in *v1.ApplyCreateCatTaskRequest) (*v1.ApplyCreateCatTaskResponse, error) {
	if in.CatId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "cat_id不能为空", errorx.ErrBadRequest)
	}
	if in.Title == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "任务标题不能为空", errorx.ErrBadRequest)
	}
	if in.ApplicantUserId == 0 {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "用户未登录", errorx.ErrUnauthorized)
	}

	apply := &task.CatTaskApply{
		CatID:           in.CatId,
		ApplicantUserID: in.ApplicantUserId,
		Title:           in.Title,
		TaskType:        in.TaskType,
		Summary:         in.Summary,
		Detail:          in.Detail,
		Status:          task.TaskApplyStatusPending,
	}

	if in.Deadline != nil {
		t := in.Deadline.AsTime()
		apply.DeadlineAt = &t
	}

	if err := l.svcCtx.Repo.Task.CreateTaskApply(l.ctx, apply); err != nil {
		logger.Error("申请创建任务失败", zap.Error(err))
		return nil, errorx.WrapDBInsert("申请创建任务失败", err)
	}

	logger.Info("志愿者申请创建任务成功", zap.Uint64("apply_id", apply.ID))

	return &v1.ApplyCreateCatTaskResponse{
		ApplyId: apply.ID,
	}, nil
}
