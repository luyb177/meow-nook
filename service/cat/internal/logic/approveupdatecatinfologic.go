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

type ApproveUpdateCatInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApproveUpdateCatInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveUpdateCatInfoLogic {
	return &ApproveUpdateCatInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApproveUpdateCatInfoLogic) ApproveUpdateCatInfo(in *v1.ApproveUpdateCatInfoRequest) (*v1.ApproveUpdateCatInfoResponse, error) {
	// TODO(permission): check that in.OperatorId has admin role

	apply, err := l.svcCtx.Repo.Cat.GetCatUpdateApplyByID(l.ctx, in.ApplyId)
	if err != nil {
		return nil, errorx.ErrCatNotFound
	}
	if apply.Status != "pending" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "申请状态不是待审核", nil)
	}

	updates := map[string]interface{}{
		"name":              apply.Name,
		"gender":            apply.Gender,
		"body_size":         apply.BodySize,
		"age_stage":         apply.AgeStage,
		"description":       apply.Description,
		"discovery_address": apply.DiscoveryAddress,
		"longitude":         apply.Longitude,
		"latitude":          apply.Latitude,
	}
	if err := l.svcCtx.Repo.Cat.UpdateCat(l.ctx, apply.CatID, updates); err != nil {
		logger.Error("ApproveUpdateCatInfo: update cat failed")
		return nil, errorx.WrapInternal("更新猫咪档案失败", err)
	}

	// Append new images from the apply
	var imageURLs []string
	if apply.ImageURLs != "" {
		_ = json.Unmarshal([]byte(apply.ImageURLs), &imageURLs)
	}
	if len(imageURLs) > 0 {
		images := make([]*cat.CatImage, 0, len(imageURLs))
		for i, u := range imageURLs {
			images = append(images, &cat.CatImage{
				CatID:      apply.CatID,
				ImageURL:   u,
				ImageType:  "normal",
				Sort:       int32(i),
				UploaderID: apply.ApplicantUserID,
			})
		}
		if err := l.svcCtx.Repo.Cat.CreateCatImages(l.ctx, images); err != nil {
			logger.Error("ApproveUpdateCatInfo: create images failed")
			// non-fatal
		}
	}

	if err := l.svcCtx.Repo.Cat.UpdateCatUpdateApplyStatus(l.ctx, apply.ID, "approved", "", in.OperatorId); err != nil {
		logger.Error("ApproveUpdateCatInfo: update apply status failed")
		return nil, errorx.WrapInternal("更新申请状态失败", err)
	}

	return &v1.ApproveUpdateCatInfoResponse{Success: true}, nil
}

// ──────────────────────────────────────────────────────────────────────────────

type RejectUpdateCatInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRejectUpdateCatInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectUpdateCatInfoLogic {
	return &RejectUpdateCatInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RejectUpdateCatInfoLogic) RejectUpdateCatInfo(in *v1.RejectUpdateCatInfoRequest) (*v1.RejectUpdateCatInfoResponse, error) {
	// TODO(permission): check that in.OperatorId has admin role

	apply, err := l.svcCtx.Repo.Cat.GetCatUpdateApplyByID(l.ctx, in.ApplyId)
	if err != nil {
		return nil, errorx.ErrCatNotFound
	}
	if apply.Status != "pending" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "申请状态不是待审核", nil)
	}

	if err := l.svcCtx.Repo.Cat.UpdateCatUpdateApplyStatus(l.ctx, apply.ID, "rejected", in.Reason, in.OperatorId); err != nil {
		logger.Error("RejectUpdateCatInfo: update apply status failed")
		return nil, errorx.WrapInternal("更新申请状态失败", err)
	}

	return &v1.RejectUpdateCatInfoResponse{Success: true}, nil
}

// ──────────────────────────────────────────────────────────────────────────────

type GetUpdateCatApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUpdateCatApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUpdateCatApplyDetailLogic {
	return &GetUpdateCatApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUpdateCatApplyDetailLogic) GetUpdateCatApplyDetail(in *v1.GetUpdateCatApplyDetailRequest) (*v1.GetUpdateCatApplyDetailResponse, error) {
	apply, err := l.svcCtx.Repo.Cat.GetCatUpdateApplyByID(l.ctx, in.ApplyId)
	if err != nil {
		return nil, errorx.ErrCatNotFound
	}

	var imageURLs []string
	if apply.ImageURLs != "" {
		_ = json.Unmarshal([]byte(apply.ImageURLs), &imageURLs)
	}

	return &v1.GetUpdateCatApplyDetailResponse{
		ApplyId:          apply.ID,
		CatId:            apply.CatID,
		Name:             apply.Name,
		Gender:           apply.Gender,
		BodySize:         apply.BodySize,
		AgeStage:         apply.AgeStage,
		Description:      apply.Description,
		DiscoveryAddress: apply.DiscoveryAddress,
		Longitude:        apply.Longitude,
		Latitude:         apply.Latitude,
		ImageUrls:        imageURLs,
		ChangeReason:     apply.ChangeReason,
		ApplicantUserId:  apply.ApplicantUserID,
		ApplicantName:    "", // TODO(user): fetch from user service
		Status:           apply.Status,
		RejectReason:     apply.RejectReason,
		CreatedAt:        timestamppb.New(apply.CreatedAt),
		UpdatedAt:        timestamppb.New(apply.UpdatedAt),
	}, nil
}

// ──────────────────────────────────────────────────────────────────────────────

type ListUpdateCatApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUpdateCatApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUpdateCatApplyLogic {
	return &ListUpdateCatApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListUpdateCatApplyLogic) ListUpdateCatApply(in *v1.ListUpdateCatApplyRequest) (*v1.ListUpdateCatApplyResponse, error) {
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	lastID := decodeCursor(in.Cursor)

	applies, hasMore, err := l.svcCtx.Repo.Cat.ListCatUpdateApplies(l.ctx, lastID, pageSize)
	if err != nil {
		logger.Error("ListUpdateCatApply: query failed")
		return nil, errorx.WrapInternal("查询申请列表失败", err)
	}

	items := make([]*v1.UpdateCatApplyItem, 0, len(applies))
	for _, a := range applies {
		items = append(items, &v1.UpdateCatApplyItem{
			ApplyId:         a.ID,
			CatId:           a.CatID,
			CatName:         "", // TODO(user): fetch cat name
			ApplicantUserId: a.ApplicantUserID,
			ApplicantName:   "", // TODO(user): fetch from user service
			Status:          a.Status,
			CreatedAt:       timestamppb.New(a.CreatedAt),
		})
	}

	nextCursor := ""
	if hasMore && len(applies) > 0 {
		nextCursor = encodeCursor(applies[len(applies)-1].ID)
	}

	return &v1.ListUpdateCatApplyResponse{
		List:       items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
