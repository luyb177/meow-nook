package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/post/internal/svc"
	postpb "github.com/luyb177/meow-nook/service/post/pb/post"
)

type CreatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CreatePostLogic) CreatePost(req *postpb.CreatePostReq) (*postpb.CreatePostResp, error) {
	if req.CatId == 0 {
		return nil, errorx.New(errorx.CodeBadRequest, "动态必须关联猫咪")
	}
	if req.Content == "" {
		return nil, errorx.New(errorx.CodeBadRequest, "动态内容不能为空")
	}
	// TODO: persist to DB
	return &postpb.CreatePostResp{PostId: 1}, nil
}

type GetPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostLogic {
	return &GetPostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetPostLogic) GetPost(req *postpb.GetPostReq) (*postpb.GetPostResp, error) {
	// TODO: fetch from DB
	return &postpb.GetPostResp{Post: &postpb.PostInfo{Id: req.PostId}}, nil
}

type DeletePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePostLogic {
	return &DeletePostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DeletePostLogic) DeletePost(req *postpb.DeletePostReq) error {
	// TODO: verify ownership; delete from DB
	return nil
}

type ListPostsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPostsLogic {
	return &ListPostsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListPostsLogic) ListPosts(req *postpb.ListPostsReq) (*postpb.ListPostsResp, error) {
	// TODO: query from DB with pagination
	return &postpb.ListPostsResp{Posts: []*postpb.PostInfo{}, Total: 0}, nil
}

type LikePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikePostLogic {
	return &LikePostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *LikePostLogic) LikePost(req *postpb.LikePostReq) error {
	// TODO: upsert like record; increment counter
	return nil
}

type UnlikePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnlikePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlikePostLogic {
	return &UnlikePostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UnlikePostLogic) UnlikePost(req *postpb.UnlikePostReq) error {
	// TODO: delete like record; decrement counter
	return nil
}

type FavoritePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoritePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoritePostLogic {
	return &FavoritePostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *FavoritePostLogic) FavoritePost(req *postpb.FavoritePostReq) error {
	// TODO: upsert favorite record; increment counter
	return nil
}

type UnfavoritePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnfavoritePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfavoritePostLogic {
	return &UnfavoritePostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UnfavoritePostLogic) UnfavoritePost(req *postpb.UnfavoritePostReq) error {
	// TODO: delete favorite record; decrement counter
	return nil
}

type AddCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCommentLogic {
	return &AddCommentLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AddCommentLogic) AddComment(req *postpb.AddCommentReq) (*postpb.AddCommentResp, error) {
	if req.Content == "" {
		return nil, errorx.New(errorx.CodeBadRequest, "评论内容不能为空")
	}
	// TODO: persist to DB; increment comment counter
	return &postpb.AddCommentResp{CommentId: 1}, nil
}

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteCommentLogic) DeleteComment(req *postpb.DeleteCommentReq) error {
	// TODO: verify ownership; delete from DB; decrement counter
	return nil
}

type ListCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentsLogic {
	return &ListCommentsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListCommentsLogic) ListComments(req *postpb.ListCommentsReq) (*postpb.ListCommentsResp, error) {
	// TODO: query from DB with pagination
	return &postpb.ListCommentsResp{Comments: []*postpb.CommentInfo{}, Total: 0}, nil
}

type ReportPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportPostLogic {
	return &ReportPostLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ReportPostLogic) ReportPost(req *postpb.ReportPostReq) error {
	if req.Reason == "" {
		return errorx.New(errorx.CodeBadRequest, "投诉原因不能为空")
	}
	// TODO: persist report; notify moderator
	return nil
}
