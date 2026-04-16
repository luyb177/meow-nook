package post

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type ListPostsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPostsLogic {
	return &ListPostsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListPostsLogic) ListPosts(req *types.ListPostsReq) (*types.ListPostsResp, error) {
	// TODO: call post gRPC service
	return &types.ListPostsResp{Posts: []types.PostInfo{}, Total: 0}, nil
}

// ──────────────────────────────────────────────

type CreatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CreatePostLogic) CreatePost(userID int64, req *types.CreatePostReq) (*types.CreatePostResp, error) {
	if req.CatId == 0 {
		return nil, errorx.New(errorx.CodeBadRequest, "动态必须关联猫咪")
	}
	if req.Content == "" {
		return nil, errorx.New(errorx.CodeBadRequest, "动态内容不能为空")
	}
	// TODO: call post gRPC service
	return &types.CreatePostResp{PostId: 1}, nil
}

// ──────────────────────────────────────────────

type GetPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostLogic {
	return &GetPostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetPostLogic) GetPost(postID int64) (*types.PostInfo, error) {
	// TODO: call post gRPC service
	return &types.PostInfo{Id: postID}, nil
}

// ──────────────────────────────────────────────

type DeletePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePostLogic {
	return &DeletePostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DeletePostLogic) DeletePost(postID, userID int64) error {
	// TODO: call post gRPC service
	return nil
}

// ──────────────────────────────────────────────

type LikePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikePostLogic {
	return &LikePostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *LikePostLogic) LikePost(postID, userID int64) error {
	// TODO: call post gRPC service
	return nil
}

// ──────────────────────────────────────────────

type UnlikePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnlikePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlikePostLogic {
	return &UnlikePostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UnlikePostLogic) UnlikePost(postID, userID int64) error {
	// TODO: call post gRPC service
	return nil
}

// ──────────────────────────────────────────────

type FavoritePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoritePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoritePostLogic {
	return &FavoritePostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *FavoritePostLogic) FavoritePost(postID, userID int64) error {
	// TODO: call post gRPC service
	return nil
}

// ──────────────────────────────────────────────

type UnfavoritePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnfavoritePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfavoritePostLogic {
	return &UnfavoritePostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UnfavoritePostLogic) UnfavoritePost(postID, userID int64) error {
	// TODO: call post gRPC service
	return nil
}

// ──────────────────────────────────────────────

type ListCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentsLogic {
	return &ListCommentsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListCommentsLogic) ListComments(postID int64, req *types.PageReq) (*types.ListCommentsResp, error) {
	// TODO: call post gRPC service
	return &types.ListCommentsResp{Comments: []types.CommentInfo{}, Total: 0}, nil
}

// ──────────────────────────────────────────────

type AddCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCommentLogic {
	return &AddCommentLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AddCommentLogic) AddComment(postID, userID int64, content string) error {
	if content == "" {
		return errorx.New(errorx.CodeBadRequest, "评论内容不能为空")
	}
	// TODO: call post gRPC service
	return nil
}

// ──────────────────────────────────────────────

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteCommentLogic) DeleteComment(commentID, userID int64) error {
	// TODO: call post gRPC service
	return nil
}

// ──────────────────────────────────────────────

type ReportPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportPostLogic {
	return &ReportPostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ReportPostLogic) ReportPost(postID, userID int64, reason string) error {
	if reason == "" {
		return errorx.New(errorx.CodeBadRequest, "投诉原因不能为空")
	}
	// TODO: call post gRPC service
	return nil
}
