package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/task"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetTaskApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskApplyDetailLogic {
	return &GetTaskApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTaskApplyDetailLogic) GetTaskApplyDetail(in *v1.GetTaskApplyDetailRequest) (*v1.GetTaskApplyDetailResponse, error) {
	logger.Info("GetTaskApplyDetail called")

	apply, err := l.svcCtx.Repo.Task.GetTaskApplyByID(l.ctx, in.ApplyId)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询任务申请失败", err)
	}

	// 权限简单校验：本人或管理员可看
	if apply.ApplicantUserID != in.RequesterId {
		// TODO: 管理员权限检查
		// casbin
	}

	vo := l.toTaskApplyVO(apply)
	return &v1.GetTaskApplyDetailResponse{Apply: vo}, nil
}

func (l *GetTaskApplyDetailLogic) toTaskApplyVO(apply *task.CatTaskApply) *v1.TaskApplyVO {
	vo := &v1.TaskApplyVO{
		Id:              apply.ID,
		CatId:           apply.CatID,
		ApplicantUserId: apply.ApplicantUserID,
		Title:           apply.Title,
		TaskType:        apply.TaskType,
		Summary:         apply.Summary,
		Detail:          apply.Detail,
		Location:        apply.Location,
		Longitude:       apply.Longitude,
		Latitude:        apply.Latitude,
		Status:          apply.Status,
		ReviewerId:      apply.ReviewerID,
		RejectReason:    apply.RejectReason,
		UrgencyLevel:    apply.UrgencyLevel,
		DifficultyLevel: apply.DifficultyLevel,
		RewardPoints:    apply.RewardPoints,
		TaskId:          apply.TaskID,
		CreatedAt:       timestamppb.New(apply.CreatedAt),
		UpdatedAt:       timestamppb.New(apply.UpdatedAt),
	}
	if apply.DeadlineAt != nil {
		vo.DeadlineAt = timestamppb.New(*apply.DeadlineAt)
	}
	if apply.ReviewedAt != nil {
		vo.ReviewedAt = timestamppb.New(*apply.ReviewedAt)
	}
	return vo
}
