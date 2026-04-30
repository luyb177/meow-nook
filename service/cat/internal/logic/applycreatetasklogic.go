package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/image"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/task"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
)

type ApplyCreateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyCreateTaskLogic {
	return &ApplyCreateTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyCreateTaskLogic) ApplyCreateTask(in *v1.ApplyCreateTaskRequest) (*v1.ApplyCreateTaskResponse, error) {
	logger.Info("ApplyCreateTask called")

	if in.CatId == 0 || in.Title == "" || in.TaskType == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "cat_id/title/task_type required", errorx.ErrBadRequest)
	}

	userID := in.ApplicantUserId

	// 防重复申请（可选）
	//exist, _ := l.svcCtx.Repo.Task.GetTaskApplyByCatAndUser(l.ctx, in.CatId, userID)
	//if exist != nil {
	//	return nil, errorx.Wrap(errorx.CodeBadRequest, "您已为该猫咪提过任务申请", errorx.ErrBadRequest)
	//}

	apply := &task.CatTaskApply{
		CatID:           in.CatId,
		ApplicantUserID: userID,
		Title:           in.Title,
		TaskType:        in.TaskType,
		Summary:         in.Summary,
		Detail:          in.Detail,
		Location:        in.Location,
		Longitude:       in.Longitude,
		Latitude:        in.Latitude,
	}

	if in.DeadlineAt != nil {
		t := in.DeadlineAt.AsTime()
		apply.DeadlineAt = &t
	}

	if err := l.svcCtx.Repo.Task.CreateTaskApply(l.ctx, apply); err != nil {
		return nil, errorx.WrapDBInsert("创建任务申请失败", err)
	}

	// 图片
	if len(in.ImageUrls) > 0 {
		imgs := make([]*image.Image, 0, len(in.ImageUrls))
		for _, url := range in.ImageUrls {
			imgs = append(imgs, &image.Image{
				TargetType: image.TargetTypeCatTaskApply,
				TargetID:   apply.ID,
				URL:        url,
				UploaderID: userID,
			})
		}
		_ = l.svcCtx.Repo.Image.BatchCreate(l.ctx, imgs)
	}

	return &v1.ApplyCreateTaskResponse{
		ApplyId: apply.ID,
		Status:  task.TaskApplyStatusPending,
		Message: "任务申请已提交",
	}, nil
}
