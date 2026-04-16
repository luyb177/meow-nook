package server

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/cat/internal/logic"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RegisterCatServer(s *grpc.Server, ctx *svc.ServiceContext) {
	catpb.RegisterCatServiceServer(s, &CatServer{ctx: ctx})
}

type CatServer struct {
	catpb.UnimplementedCatServiceServer
	ctx *svc.ServiceContext
}

func toGRPCError(err error) error {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*errorx.AppError); ok {
		switch ae.Code {
		case errorx.CodeNotFound, errorx.CodeCatNotFound:
			return status.Error(codes.NotFound, ae.Msg)
		case errorx.CodeForbidden, errorx.CodePermissionDenied:
			return status.Error(codes.PermissionDenied, ae.Msg)
		case errorx.CodeBadRequest:
			return status.Error(codes.InvalidArgument, ae.Msg)
		default:
			return status.Error(codes.Internal, ae.Msg)
		}
	}
	return status.Error(codes.Internal, "服务器内部错误")
}

func (s *CatServer) CreateCat(ctx context.Context, req *catpb.CreateCatReq) (*catpb.CreateCatResp, error) {
	resp, err := logic.NewCreateCatLogic(ctx, s.ctx).CreateCat(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *CatServer) GetCat(ctx context.Context, req *catpb.GetCatReq) (*catpb.GetCatResp, error) {
	resp, err := logic.NewGetCatLogic(ctx, s.ctx).GetCat(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *CatServer) UpdateCat(ctx context.Context, req *catpb.UpdateCatReq) (*catpb.UpdateCatResp, error) {
	if err := logic.NewUpdateCatLogic(ctx, s.ctx).UpdateCat(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &catpb.UpdateCatResp{}, nil
}

func (s *CatServer) DeleteCat(ctx context.Context, req *catpb.DeleteCatReq) (*catpb.DeleteCatResp, error) {
	if err := logic.NewDeleteCatLogic(ctx, s.ctx).DeleteCat(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &catpb.DeleteCatResp{}, nil
}

func (s *CatServer) ListCats(ctx context.Context, req *catpb.ListCatsReq) (*catpb.ListCatsResp, error) {
	resp, err := logic.NewListCatsLogic(ctx, s.ctx).ListCats(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *CatServer) AddRescueRecord(ctx context.Context, req *catpb.AddRescueRecordReq) (*catpb.AddRescueRecordResp, error) {
	resp, err := logic.NewAddRescueRecordLogic(ctx, s.ctx).AddRescueRecord(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *CatServer) ListRescueRecords(ctx context.Context, req *catpb.ListRescueRecordsReq) (*catpb.ListRescueRecordsResp, error) {
	resp, err := logic.NewListRescueRecordsLogic(ctx, s.ctx).ListRescueRecords(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *CatServer) AddHealthRecord(ctx context.Context, req *catpb.AddHealthRecordReq) (*catpb.AddHealthRecordResp, error) {
	resp, err := logic.NewAddHealthRecordLogic(ctx, s.ctx).AddHealthRecord(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *CatServer) ListHealthRecords(ctx context.Context, req *catpb.ListHealthRecordsReq) (*catpb.ListHealthRecordsResp, error) {
	resp, err := logic.NewListHealthRecordsLogic(ctx, s.ctx).ListHealthRecords(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *CatServer) GetStats(ctx context.Context, req *catpb.GetStatsReq) (*catpb.GetStatsResp, error) {
	resp, err := logic.NewGetStatsLogic(ctx, s.ctx).GetStats(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *CatServer) GetHeatmap(ctx context.Context, req *catpb.GetHeatmapReq) (*catpb.GetHeatmapResp, error) {
	resp, err := logic.NewGetHeatmapLogic(ctx, s.ctx).GetHeatmap(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}
