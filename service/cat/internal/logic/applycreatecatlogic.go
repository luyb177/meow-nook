package logic

import (
	"context"
	"encoding/json"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyCreateCatLogic {
	return &ApplyCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 申请创建小猫档案
func (l *ApplyCreateCatLogic) ApplyCreateCat(in *v1.ApplyCreateCatRequest) (*v1.ApplyCreateCatResponse, error) {
	imageURLsJSON, err := json.Marshal(in.ImageUrls)
	if err != nil {
		logger.Error("ApplyCreateCat: marshal image_urls failed")
		return nil, errorx.WrapInternal("序列化图片URL失败", err)
	}

	apply := &cat.CatCreateApply{
		Name:             in.Name,
		Gender:           in.Gender,
		BodySize:         in.BodySize,
		AgeStage:         in.AgeStage,
		Description:      in.Description,
		DiscoveryAddress: in.DiscoveryAddress,
		Longitude:        in.Longitude,
		Latitude:         in.Latitude,
		ImageURLs:        string(imageURLsJSON),
		ApplicantUserID:  in.ApplicantUserId,
		Status:           "pending",
	}

	if err := l.svcCtx.Repo.Cat.CreateCatCreateApply(l.ctx, apply); err != nil {
		logger.Error("ApplyCreateCat: create apply failed")
		return nil, errorx.WrapInternal("创建申请失败", err)
	}

	return &v1.ApplyCreateCatResponse{ApplyId: apply.ID}, nil
}

// ──────────────────────────────────────────────────────────────────────────────

type GetMyCreateCatApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMyCreateCatApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyCreateCatApplyDetailLogic {
	return &GetMyCreateCatApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查看自己提交的创建申请详情
func (l *GetMyCreateCatApplyDetailLogic) GetMyCreateCatApplyDetail(in *v1.GetMyCreateCatApplyDetailRequest) (*v1.GetMyCreateCatApplyDetailResponse, error) {
	apply, err := l.svcCtx.Repo.Cat.GetCatCreateApplyByIDAndUser(l.ctx, in.ApplyId, in.UserId)
	if err != nil {
		return nil, errorx.ErrCatNotFound
	}

	return &v1.GetMyCreateCatApplyDetailResponse{
		ApplyId:      apply.ID,
		Status:       apply.Status,
		RejectReason: apply.RejectReason,
		CreatedAt:    timestamppb.New(apply.CreatedAt),
		UpdatedAt:    timestamppb.New(apply.UpdatedAt),
	}, nil
}
