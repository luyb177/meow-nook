package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

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
	// TODO(permission): check that in.OperatorId has admin role

	apply, err := l.svcCtx.Repo.Cat.GetCatTaskApplyByID(l.ctx, in.ApplyId)
	if err != nil {
		return nil, errorx.ErrTaskNotFound
	}
	if apply.Status != "pending" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "申请状态不是待审核", nil)
	}

	task := &cat.CatTask{
		CatID:           apply.CatID,
		Title:           apply.Title,
		TaskType:        apply.TaskType,
		UrgencyLevel:    in.UrgencyLevel,
		DifficultyLevel: in.DifficultyLevel,
		RewardPoints:    in.RewardPoints,
		Status:          "pending",
		Summary:         apply.Summary,
		Detail:          apply.Detail,
		Deadline:        apply.Deadline,
	}
	if err := l.svcCtx.Repo.Cat.CreateCatTask(l.ctx, task); err != nil {
		logger.Error("ApproveCreateCatTask: create task failed")
		return nil, errorx.WrapInternal("创建任务失败", err)
	}

	if err := l.svcCtx.Repo.Cat.UpdateCatTaskApply(l.ctx, apply.ID, map[string]interface{}{
		"status":           "approved",
		"reviewer_id":      in.OperatorId,
		"urgency_level":    in.UrgencyLevel,
		"difficulty_level": in.DifficultyLevel,
		"reward_points":    in.RewardPoints,
	}); err != nil {
		logger.Error("ApproveCreateCatTask: update apply status failed")
		return nil, errorx.WrapInternal("更新申请状态失败", err)
	}

	return &v1.ApproveCreateCatTaskResponse{TaskId: task.ID}, nil
}

// ──────────────────────────────────────────────────────────────────────────────

type RejectCreateCatTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRejectCreateCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectCreateCatTaskLogic {
	return &RejectCreateCatTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RejectCreateCatTaskLogic) RejectCreateCatTask(in *v1.RejectCreateCatTaskRequest) (*v1.RejectCreateCatTaskResponse, error) {
	// TODO(permission): check that in.OperatorId has admin role

	apply, err := l.svcCtx.Repo.Cat.GetCatTaskApplyByID(l.ctx, in.ApplyId)
	if err != nil {
		return nil, errorx.ErrTaskNotFound
	}
	if apply.Status != "pending" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "申请状态不是待审核", nil)
	}

	if err := l.svcCtx.Repo.Cat.UpdateCatTaskApply(l.ctx, apply.ID, map[string]interface{}{
		"status":        "rejected",
		"reject_reason": in.Reason,
		"reviewer_id":   in.OperatorId,
	}); err != nil {
		logger.Error("RejectCreateCatTask: update apply status failed")
		return nil, errorx.WrapInternal("更新申请状态失败", err)
	}

	return &v1.RejectCreateCatTaskResponse{Success: true}, nil
}

// ──────────────────────────────────────────────────────────────────────────────

type GetCreateCatTaskApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCreateCatTaskApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCreateCatTaskApplyDetailLogic {
	return &GetCreateCatTaskApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCreateCatTaskApplyDetailLogic) GetCreateCatTaskApplyDetail(in *v1.GetCreateCatTaskApplyDetailRequest) (*v1.GetCreateCatTaskApplyDetailResponse, error) {
	apply, err := l.svcCtx.Repo.Cat.GetCatTaskApplyByID(l.ctx, in.ApplyId)
	if err != nil {
		return nil, errorx.ErrTaskNotFound
	}

	return &v1.GetCreateCatTaskApplyDetailResponse{
		ApplyId:         apply.ID,
		CatId:           apply.CatID,
		Title:           apply.Title,
		TaskType:        apply.TaskType,
		Summary:         apply.Summary,
		Detail:          apply.Detail,
		Deadline:        timestamppb.New(apply.Deadline),
		ApplicantUserId: apply.ApplicantUserID,
		ApplicantName:   "", // TODO(user): fetch from user service
		Status:          apply.Status,
		RejectReason:    apply.RejectReason,
		CreatedAt:       timestamppb.New(apply.CreatedAt),
		UpdatedAt:       timestamppb.New(apply.UpdatedAt),
	}, nil
}

// ──────────────────────────────────────────────────────────────────────────────

type ListCreateCatTaskApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCreateCatTaskApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCreateCatTaskApplyLogic {
	return &ListCreateCatTaskApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCreateCatTaskApplyLogic) ListCreateCatTaskApply(in *v1.ListCreateCatTaskApplyRequest) (*v1.ListCreateCatTaskApplyResponse, error) {
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	lastID := decodeCursor(in.Cursor)

	applies, hasMore, err := l.svcCtx.Repo.Cat.ListCatTaskApplies(l.ctx, lastID, pageSize)
	if err != nil {
		logger.Error("ListCreateCatTaskApply: query failed")
		return nil, errorx.WrapInternal("查询任务申请列表失败", err)
	}

	items := make([]*v1.CreateCatTaskApplyItem, 0, len(applies))
	for _, a := range applies {
		items = append(items, &v1.CreateCatTaskApplyItem{
			ApplyId:         a.ID,
			CatId:           a.CatID,
			Title:           a.Title,
			TaskType:        a.TaskType,
			ApplicantUserId: a.ApplicantUserID,
			ApplicantName:   "", // TODO(user): fetch from user service
			Status:          a.Status,
			CreatedAt:       timestamppb.New(a.CreatedAt),
		})
	}

	nextCursor := ""
	if hasMore && len(applies) > 0 {
		nextCursor = encodeCursor(applies[len(applies)-1].ID)
	}

	return &v1.ListCreateCatTaskApplyResponse{
		List:       items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
