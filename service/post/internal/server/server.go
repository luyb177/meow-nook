package server

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/post/internal/logic"
	"github.com/luyb177/meow-nook/service/post/internal/svc"
	postpb "github.com/luyb177/meow-nook/service/post/pb/post"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RegisterPostServer(s *grpc.Server, ctx *svc.ServiceContext) {
	postpb.RegisterPostServiceServer(s, &PostServer{ctx: ctx})
}

type PostServer struct {
	postpb.UnimplementedPostServiceServer
	ctx *svc.ServiceContext
}

func toGRPCError(err error) error {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*errorx.AppError); ok {
		switch ae.Code {
		case errorx.CodeNotFound, errorx.CodePostNotFound:
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

func (s *PostServer) CreatePost(ctx context.Context, req *postpb.CreatePostReq) (*postpb.CreatePostResp, error) {
	resp, err := logic.NewCreatePostLogic(ctx, s.ctx).CreatePost(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *PostServer) GetPost(ctx context.Context, req *postpb.GetPostReq) (*postpb.GetPostResp, error) {
	resp, err := logic.NewGetPostLogic(ctx, s.ctx).GetPost(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *PostServer) DeletePost(ctx context.Context, req *postpb.DeletePostReq) (*postpb.DeletePostResp, error) {
	if err := logic.NewDeletePostLogic(ctx, s.ctx).DeletePost(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &postpb.DeletePostResp{}, nil
}

func (s *PostServer) ListPosts(ctx context.Context, req *postpb.ListPostsReq) (*postpb.ListPostsResp, error) {
	resp, err := logic.NewListPostsLogic(ctx, s.ctx).ListPosts(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *PostServer) LikePost(ctx context.Context, req *postpb.LikePostReq) (*postpb.LikePostResp, error) {
	if err := logic.NewLikePostLogic(ctx, s.ctx).LikePost(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &postpb.LikePostResp{}, nil
}

func (s *PostServer) UnlikePost(ctx context.Context, req *postpb.UnlikePostReq) (*postpb.UnlikePostResp, error) {
	if err := logic.NewUnlikePostLogic(ctx, s.ctx).UnlikePost(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &postpb.UnlikePostResp{}, nil
}

func (s *PostServer) FavoritePost(ctx context.Context, req *postpb.FavoritePostReq) (*postpb.FavoritePostResp, error) {
	if err := logic.NewFavoritePostLogic(ctx, s.ctx).FavoritePost(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &postpb.FavoritePostResp{}, nil
}

func (s *PostServer) UnfavoritePost(ctx context.Context, req *postpb.UnfavoritePostReq) (*postpb.UnfavoritePostResp, error) {
	if err := logic.NewUnfavoritePostLogic(ctx, s.ctx).UnfavoritePost(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &postpb.UnfavoritePostResp{}, nil
}

func (s *PostServer) AddComment(ctx context.Context, req *postpb.AddCommentReq) (*postpb.AddCommentResp, error) {
	resp, err := logic.NewAddCommentLogic(ctx, s.ctx).AddComment(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *PostServer) DeleteComment(ctx context.Context, req *postpb.DeleteCommentReq) (*postpb.DeleteCommentResp, error) {
	if err := logic.NewDeleteCommentLogic(ctx, s.ctx).DeleteComment(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &postpb.DeleteCommentResp{}, nil
}

func (s *PostServer) ListComments(ctx context.Context, req *postpb.ListCommentsReq) (*postpb.ListCommentsResp, error) {
	resp, err := logic.NewListCommentsLogic(ctx, s.ctx).ListComments(req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (s *PostServer) ReportPost(ctx context.Context, req *postpb.ReportPostReq) (*postpb.ReportPostResp, error) {
	if err := logic.NewReportPostLogic(ctx, s.ctx).ReportPost(req); err != nil {
		return nil, toGRPCError(err)
	}
	return &postpb.ReportPostResp{}, nil
}
