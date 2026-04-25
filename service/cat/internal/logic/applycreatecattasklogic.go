package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

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
	if _, err := l.svcCtx.Repo.Cat.GetCatByID(l.ctx, in.CatId); err != nil {
		return nil, errorx.ErrCatNotFound
	}

	apply := &cat.CatTaskApply{
		CatID:           in.CatId,
		Title:           in.Title,
		TaskType:        in.TaskType,
		Summary:         in.Summary,
		Detail:          in.Detail,
		Deadline:        in.Deadline.AsTime(),
		ApplicantUserID: in.ApplicantUserId,
		Status:          "pending",
	}

	if err := l.svcCtx.Repo.Cat.CreateCatTaskApply(l.ctx, apply); err != nil {
		logger.Error("ApplyCreateCatTask: create apply failed")
		return nil, errorx.WrapInternal("创建任务申请失败", err)
	}

	return &v1.ApplyCreateCatTaskResponse{ApplyId: apply.ID}, nil
}
