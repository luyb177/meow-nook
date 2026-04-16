package server

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/task/internal/logic"
	"github.com/luyb177/meow-nook/service/task/internal/svc"
	taskpb "github.com/luyb177/meow-nook/service/task/pb/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RegisterTaskServer(s *grpc.Server, ctx *svc.ServiceContext) {
	taskpb.RegisterTaskServiceServer(s, &TaskServer{ctx: ctx})
}

type TaskServer struct {
	taskpb.UnimplementedTaskServiceServer
	ctx *svc.ServiceContext
}

func toGRPCError(err error) error {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*errorx.AppError); ok {
		switch ae.Code {
		case errorx.CodeNotFound, errorx.CodeTaskNotFound:
			return status.Error(codes.NotFound, ae.Msg)
		case errorx.CodeConflict, errorx.CodeTaskAlreadyClaimed, errorx.CodeTaskFull:
			return status.Error(codes.AlreadyExists, ae.Msg)
		case errorx.CodeForbidden, errorx.CodePermissionDenied, errorx.CodeTaskNotOwned:
			return status.Error(codes.PermissionDenied, ae.Msg)
		case errorx.CodeBadRequest:
			return status.Error(codes.InvalidArgument, ae.Msg)
		default:
			return status.Error(codes.Internal, ae.Msg)
		}
	}
	return status.Error(codes.Internal, "服务器内部错误")
}

func (s *TaskServer) CreateTask(ctx context.Context, req *taskpb.CreateTaskReq) (*taskpb.CreateTaskResp, error) {
	resp, err := logic.NewCreateTaskLogic(ctx, s.ctx).CreateTask(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *TaskServer) GetTask(ctx context.Context, req *taskpb.GetTaskReq) (*taskpb.GetTaskResp, error) {
	resp, err := logic.NewGetTaskLogic(ctx, s.ctx).GetTask(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *TaskServer) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskReq) (*taskpb.UpdateTaskResp, error) {
	if err := logic.NewUpdateTaskLogic(ctx, s.ctx).UpdateTask(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &taskpb.UpdateTaskResp{}, nil
}

func (s *TaskServer) ListTasks(ctx context.Context, req *taskpb.ListTasksReq) (*taskpb.ListTasksResp, error) {
	resp, err := logic.NewListTasksLogic(ctx, s.ctx).ListTasks(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *TaskServer) ClaimTask(ctx context.Context, req *taskpb.ClaimTaskReq) (*taskpb.ClaimTaskResp, error) {
	if err := logic.NewClaimTaskLogic(ctx, s.ctx).ClaimTask(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &taskpb.ClaimTaskResp{}, nil
}

func (s *TaskServer) AbandonTask(ctx context.Context, req *taskpb.AbandonTaskReq) (*taskpb.AbandonTaskResp, error) {
	if err := logic.NewAbandonTaskLogic(ctx, s.ctx).AbandonTask(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &taskpb.AbandonTaskResp{}, nil
}

func (s *TaskServer) CompleteTask(ctx context.Context, req *taskpb.CompleteTaskReq) (*taskpb.CompleteTaskResp, error) {
	if err := logic.NewCompleteTaskLogic(ctx, s.ctx).CompleteTask(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &taskpb.CompleteTaskResp{}, nil
}

func (s *TaskServer) ExtendDeadline(ctx context.Context, req *taskpb.ExtendDeadlineReq) (*taskpb.ExtendDeadlineResp, error) {
	if err := logic.NewExtendDeadlineLogic(ctx, s.ctx).ExtendDeadline(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &taskpb.ExtendDeadlineResp{}, nil
}
