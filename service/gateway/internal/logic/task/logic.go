package task

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type ListTasksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTasksLogic {
	return &ListTasksLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListTasksLogic) ListTasks(req *types.ListTasksReq) (*types.ListTasksResp, error) {
	// TODO: call task gRPC service
	return &types.ListTasksResp{Tasks: []types.TaskInfo{}, Total: 0}, nil
}

// ──────────────────────────────────────────────

type CreateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskLogic {
	return &CreateTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateTaskLogic) CreateTask(userID int64, req *types.CreateTaskReq) (*types.CreateTaskResp, error) {
	if req.Title == "" {
		return nil, errorx.New(errorx.CodeBadRequest, "任务标题不能为空")
	}
	if req.CatId == 0 {
		return nil, errorx.New(errorx.CodeBadRequest, "必须关联猫咪")
	}
	// TODO: call task gRPC service
	return &types.CreateTaskResp{TaskId: 1}, nil
}

// ──────────────────────────────────────────────

type GetTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskLogic {
	return &GetTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetTaskLogic) GetTask(taskID int64) (*types.GetTaskResp, error) {
	// TODO: call task gRPC service
	return &types.GetTaskResp{
		Task:      types.TaskInfo{Id: taskID},
		Logs:      []types.TaskLog{},
		Assignees: []types.TaskAssignee{},
	}, nil
}

// ──────────────────────────────────────────────

type UpdateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskLogic {
	return &UpdateTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateTaskLogic) UpdateTask(taskID, userID int64, req *types.UpdateTaskReq) error {
	// TODO: call task gRPC service
	return nil
}

// ──────────────────────────────────────────────

type ClaimTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClaimTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimTaskLogic {
	return &ClaimTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ClaimTaskLogic) ClaimTask(taskID, userID int64) error {
	// TODO: call task gRPC service
	return nil
}

// ──────────────────────────────────────────────

type AbandonTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAbandonTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AbandonTaskLogic {
	return &AbandonTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AbandonTaskLogic) AbandonTask(taskID, userID int64, reason string) error {
	// TODO: call task gRPC service
	return nil
}

// ──────────────────────────────────────────────

type CompleteTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompleteTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteTaskLogic {
	return &CompleteTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CompleteTaskLogic) CompleteTask(taskID, userID int64, proof string) error {
	// TODO: call task gRPC service
	return nil
}

// ──────────────────────────────────────────────

type ExtendDeadlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExtendDeadlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExtendDeadlineLogic {
	return &ExtendDeadlineLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ExtendDeadlineLogic) ExtendDeadline(taskID, userID, newDeadline int64) error {
	if newDeadline == 0 {
		return errorx.New(errorx.CodeBadRequest, "截止时间不能为空")
	}
	// TODO: call task gRPC service
	return nil
}
