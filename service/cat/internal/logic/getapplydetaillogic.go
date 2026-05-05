// internal/logic/get_apply_detail_logic.go
package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/image"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/tag"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"go.uber.org/zap"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApplyDetailLogic {
	return &GetApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetApplyDetailLogic) GetApplyDetail(in *v1.GetApplyDetailRequest) (*v1.GetApplyDetailResponse, error) {
	// 1. 查询申请
	apply, err := l.svcCtx.Repo.Cat.GetApplyByID(l.ctx, in.ApplyId)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询申请失败", err)
	}

	// 2. 查询图片
	imgs, err := l.svcCtx.Repo.Image.ListByTarget(
		l.ctx,
		image.TargetTypeCatApply,
		apply.ID,
	)
	if err != nil {
		logger.Error("查询申请图片失败", zap.Error(err))
	}

	images := make([]*v1.ImageItem, 0, len(imgs))
	for _, img := range imgs {
		images = append(images, &v1.ImageItem{
			Url:         img.URL,
			Description: img.Description,
			Sort:        img.Sort,
			IsCover:     img.IsCover,
		})
	}

	// 3. 查询标签（可选）
	tags, err := l.svcCtx.Repo.Tag.GetTagsByTarget(
		l.ctx,
		tag.TargetTypeCatApply,
		apply.ID,
	)
	if err != nil {
		logger.Error("查询申请标签失败", zap.Error(err))
	}

	_ = tags // 如果需要返回标签，可以加到 response 里

	// 4. 组装返回
	resp := &v1.GetApplyDetailResponse{
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
		ApplicantUserId:  apply.ApplicantUserID,
		Status:           apply.Status,
		RejectReason:     apply.RejectReason,
		CancelReason:     apply.CancelReason,
		ReviewerId:       apply.ReviewerID,
		Images:           images,
		CreatedAt:        timestamppb.New(apply.CreatedAt),
		UpdatedAt:        timestamppb.New(apply.UpdatedAt),
	}

	return resp, nil
}
