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

type ApplyUpdateCatInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyUpdateCatInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyUpdateCatInfoLogic {
	return &ApplyUpdateCatInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 申请修改猫咪档案信息
func (l *ApplyUpdateCatInfoLogic) ApplyUpdateCatInfo(in *v1.ApplyUpdateCatInfoRequest) (*v1.ApplyUpdateCatInfoResponse, error) {
	if _, err := l.svcCtx.Repo.Cat.GetCatByID(l.ctx, in.CatId); err != nil {
		return nil, errorx.ErrCatNotFound
	}

	imageURLsJSON, err := json.Marshal(in.ImageUrls)
	if err != nil {
		logger.Error("ApplyUpdateCatInfo: marshal image_urls failed")
		return nil, errorx.WrapInternal("序列化图片URL失败", err)
	}

	apply := &cat.CatUpdateApply{
		CatID:            in.CatId,
		Name:             in.Name,
		Gender:           in.Gender,
		BodySize:         in.BodySize,
		AgeStage:         in.AgeStage,
		Description:      in.Description,
		DiscoveryAddress: in.DiscoveryAddress,
		Longitude:        in.Longitude,
		Latitude:         in.Latitude,
		ImageURLs:        string(imageURLsJSON),
		ChangeReason:     in.ChangeReason,
		ApplicantUserID:  in.ApplicantUserId,
		Status:           "pending",
	}

	if err := l.svcCtx.Repo.Cat.CreateCatUpdateApply(l.ctx, apply); err != nil {
		logger.Error("ApplyUpdateCatInfo: create apply failed")
		return nil, errorx.WrapInternal("创建修改申请失败", err)
	}

	return &v1.ApplyUpdateCatInfoResponse{ApplyId: apply.ID}, nil
}

// ──────────────────────────────────────────────────────────────────────────────

type GetMyUpdateCatApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMyUpdateCatApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyUpdateCatApplyDetailLogic {
	return &GetMyUpdateCatApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查看自己提交的修改申请详情
func (l *GetMyUpdateCatApplyDetailLogic) GetMyUpdateCatApplyDetail(in *v1.GetMyUpdateCatApplyDetailRequest) (*v1.GetMyUpdateCatApplyDetailResponse, error) {
	apply, err := l.svcCtx.Repo.Cat.GetCatUpdateApplyByIDAndUser(l.ctx, in.ApplyId, in.UserId)
	if err != nil {
		return nil, errorx.ErrCatNotFound
	}

	return &v1.GetMyUpdateCatApplyDetailResponse{
		ApplyId:      apply.ID,
		Status:       apply.Status,
		RejectReason: apply.RejectReason,
		CreatedAt:    timestamppb.New(apply.CreatedAt),
		UpdatedAt:    timestamppb.New(apply.UpdatedAt),
	}, nil
}
