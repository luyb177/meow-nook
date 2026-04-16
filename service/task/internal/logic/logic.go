package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/task/internal/svc"
	taskpb "github.com/luyb177/meow-nook/service/task/pb/task"
)

type CreateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskLogic {
	return &CreateTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateTaskLogic) CreateTask(req *taskpb.CreateTaskReq) (*taskpb.CreateTaskResp, error) {
	if req.Title == "" {
		return nil, errorx.New(errorx.CodeBadRequest, "任务标题不能为空")
	}
	// TODO: persist to DB
	return &taskpb.CreateTaskResp{TaskId: 1}, nil
}

type GetTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskLogic {
	return &GetTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetTaskLogic) GetTask(req *taskpb.GetTaskReq) (*taskpb.GetTaskResp, error) {
	// TODO: fetch from DB
	return &taskpb.GetTaskResp{
		Task:      &taskpb.TaskInfo{Id: req.TaskId},
		Logs:      []*taskpb.TaskLog{},
		Assignees: []*taskpb.TaskAssignee{},
	}, nil
}

type UpdateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskLogic {
	return &UpdateTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateTaskLogic) UpdateTask(req *taskpb.UpdateTaskReq) error {
	// TODO: update in DB; record log
	return nil
}

type ListTasksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTasksLogic {
	return &ListTasksLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListTasksLogic) ListTasks(req *taskpb.ListTasksReq) (*taskpb.ListTasksResp, error) {
	// TODO: query from DB with geo filter support
	return &taskpb.ListTasksResp{Tasks: []*taskpb.TaskInfo{}, Total: 0}, nil
}

type ClaimTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClaimTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimTaskLogic {
	return &ClaimTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ClaimTaskLogic) ClaimTask(req *taskpb.ClaimTaskReq) error {
	// TODO: check capacity; add assignee; update status
	return nil
}

type AbandonTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAbandonTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AbandonTaskLogic {
	return &AbandonTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AbandonTaskLogic) AbandonTask(req *taskpb.AbandonTaskReq) error {
	// TODO: update assignee status; deduct trust points
	return nil
}

type CompleteTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompleteTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteTaskLogic {
	return &CompleteTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CompleteTaskLogic) CompleteTask(req *taskpb.CompleteTaskReq) error {
	// TODO: mark as completed; award points to assignee
	return nil
}

type ExtendDeadlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExtendDeadlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExtendDeadlineLogic {
	return &ExtendDeadlineLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ExtendDeadlineLogic) ExtendDeadline(req *taskpb.ExtendDeadlineReq) error {
	if req.NewDeadline == 0 {
		return errorx.New(errorx.CodeBadRequest, "截止时间不能为空")
	}
	// TODO: update in DB; record log
	return nil
}
