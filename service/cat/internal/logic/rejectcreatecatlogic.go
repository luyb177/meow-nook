package logic

import (
	"context"
	"encoding/json"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RejectCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRejectCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectCreateCatLogic {
	return &RejectCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RejectCreateCatLogic) RejectCreateCat(in *v1.RejectCreateCatRequest) (*v1.RejectCreateCatResponse, error) {
	// TODO(permission): check that in.OperatorId has admin role

	apply, err := l.svcCtx.Repo.Cat.GetCatCreateApplyByID(l.ctx, in.ApplyId)
	if err != nil {
		return nil, errorx.ErrCatNotFound
	}
	if apply.Status != "pending" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "申请状态不是待审核", nil)
	}

	if err := l.svcCtx.Repo.Cat.UpdateCatCreateApplyStatus(l.ctx, apply.ID, "rejected", in.Reason, in.OperatorId); err != nil {
		logger.Error("RejectCreateCat: update apply status failed")
		return nil, errorx.WrapInternal("更新申请状态失败", err)
	}

	return &v1.RejectCreateCatResponse{Success: true}, nil
}

// ──────────────────────────────────────────────────────────────────────────────

type GetCreateCatApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCreateCatApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCreateCatApplyDetailLogic {
	return &GetCreateCatApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCreateCatApplyDetailLogic) GetCreateCatApplyDetail(in *v1.GetCreateCatApplyDetailRequest) (*v1.GetCreateCatApplyDetailResponse, error) {
	apply, err := l.svcCtx.Repo.Cat.GetCatCreateApplyByID(l.ctx, in.ApplyId)
	if err != nil {
		return nil, errorx.ErrCatNotFound
	}

	var imageURLs []string
	if apply.ImageURLs != "" {
		_ = json.Unmarshal([]byte(apply.ImageURLs), &imageURLs)
	}

	return &v1.GetCreateCatApplyDetailResponse{
		ApplyId:         apply.ID,
		Name:            apply.Name,
		Gender:          apply.Gender,
		BodySize:        apply.BodySize,
		AgeStage:        apply.AgeStage,
		Description:     apply.Description,
		DiscoveryAddress: apply.DiscoveryAddress,
		Longitude:       apply.Longitude,
		Latitude:        apply.Latitude,
		ImageUrls:       imageURLs,
		ApplicantUserId: apply.ApplicantUserID,
		ApplicantName:   "", // TODO(user): fetch from user service
		Status:          apply.Status,
		RejectReason:    apply.RejectReason,
		CreatedAt:       timestamppb.New(apply.CreatedAt),
		UpdatedAt:       timestamppb.New(apply.UpdatedAt),
	}, nil
}

// ──────────────────────────────────────────────────────────────────────────────

type ListCreateCatApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCreateCatApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCreateCatApplyLogic {
	return &ListCreateCatApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCreateCatApplyLogic) ListCreateCatApply(in *v1.ListCreateCatApplyRequest) (*v1.ListCreateCatApplyResponse, error) {
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	lastID := decodeCursor(in.Cursor)

	applies, hasMore, err := l.svcCtx.Repo.Cat.ListCatCreateApplies(l.ctx, lastID, pageSize)
	if err != nil {
		logger.Error("ListCreateCatApply: query failed")
		return nil, errorx.WrapInternal("查询申请列表失败", err)
	}

	items := make([]*v1.CreateCatApplyItem, 0, len(applies))
	for _, a := range applies {
		items = append(items, &v1.CreateCatApplyItem{
			ApplyId:         a.ID,
			Name:            a.Name,
			Gender:          a.Gender,
			AgeStage:        a.AgeStage,
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

	return &v1.ListCreateCatApplyResponse{
		List:       items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
