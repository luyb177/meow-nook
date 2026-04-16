package server

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/adoption/internal/logic"
	"github.com/luyb177/meow-nook/service/adoption/internal/svc"
	adoptionpb "github.com/luyb177/meow-nook/service/adoption/pb/adoption"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RegisterAdoptionServer(s *grpc.Server, ctx *svc.ServiceContext) {
	adoptionpb.RegisterAdoptionServiceServer(s, &AdoptionServer{ctx: ctx})
}

type AdoptionServer struct {
	adoptionpb.UnimplementedAdoptionServiceServer
	ctx *svc.ServiceContext
}

func toGRPCError(err error) error {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*errorx.AppError); ok {
		switch ae.Code {
		case errorx.CodeNotFound, errorx.CodeAdoptionNotFound:
			return status.Error(codes.NotFound, ae.Msg)
		case errorx.CodeConflict, errorx.CodeAdoptionAlreadyApplied:
			return status.Error(codes.AlreadyExists, ae.Msg)
		case errorx.CodeForbidden, errorx.CodeInsufficientPoints:
			return status.Error(codes.PermissionDenied, ae.Msg)
		case errorx.CodeBadRequest:
			return status.Error(codes.InvalidArgument, ae.Msg)
		default:
			return status.Error(codes.Internal, ae.Msg)
		}
	}
	return status.Error(codes.Internal, "服务器内部错误")
}

func (s *AdoptionServer) ApplyAdoption(ctx context.Context, req *adoptionpb.ApplyAdoptionReq) (*adoptionpb.ApplyAdoptionResp, error) {
	resp, err := logic.NewApplyAdoptionLogic(ctx, s.ctx).ApplyAdoption(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *AdoptionServer) ReviewApplication(ctx context.Context, req *adoptionpb.ReviewApplicationReq) (*adoptionpb.ReviewApplicationResp, error) {
	if err := logic.NewReviewApplicationLogic(ctx, s.ctx).ReviewApplication(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &adoptionpb.ReviewApplicationResp{}, nil
}

func (s *AdoptionServer) GetApplication(ctx context.Context, req *adoptionpb.GetApplicationReq) (*adoptionpb.GetApplicationResp, error) {
	resp, err := logic.NewGetApplicationLogic(ctx, s.ctx).GetApplication(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *AdoptionServer) ListApplications(ctx context.Context, req *adoptionpb.ListApplicationsReq) (*adoptionpb.ListApplicationsResp, error) {
	resp, err := logic.NewListApplicationsLogic(ctx, s.ctx).ListApplications(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *AdoptionServer) SubmitFollowUp(ctx context.Context, req *adoptionpb.SubmitFollowUpReq) (*adoptionpb.SubmitFollowUpResp, error) {
	resp, err := logic.NewSubmitFollowUpLogic(ctx, s.ctx).SubmitFollowUp(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *AdoptionServer) ListFollowUps(ctx context.Context, req *adoptionpb.ListFollowUpsReq) (*adoptionpb.ListFollowUpsResp, error) {
	resp, err := logic.NewListFollowUpsLogic(ctx, s.ctx).ListFollowUps(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}
