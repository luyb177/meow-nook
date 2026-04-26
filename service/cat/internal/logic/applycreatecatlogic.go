package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/image"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/tag"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ApplyCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyCreateCatLogic {
	return &ApplyCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyCreateCatLogic) ApplyCreateCat(in *v1.ApplyCreateCatRequest) (*v1.ApplyCreateCatResponse, error) {
	var applyID uint64

	err := l.svcCtx.Repo.WithTx(l.ctx, func(tx *gorm.DB) error {
		// ========== 1. 创建申请记录 ==========
		apply := &cat.CatCreateApply{
			Name:             in.GetName(),
			Gender:           in.GetGender(),
			BodySize:         in.GetBodySize(),
			AgeStage:         in.GetAgeStage(),
			Description:      in.GetDescription(),
			DiscoveryAddress: in.GetDiscoveryAddress(),
			Longitude:        in.GetLongitude(),
			Latitude:         in.GetLatitude(),
			ApplicantUserID:  in.GetApplicantUserId(),
			Status:           cat.ApplyStatusPending,
		}

		if err := l.svcCtx.Repo.Cat.CreateApply(l.ctx, apply, tx); err != nil {
			logger.Error("申请创建猫咪档案失败", zap.Error(err))
			return errorx.WrapDBInsert("申请创建猫咪档案失败", err)
		}

		applyID = apply.ID

		// ========== 2. 创建图片记录 ==========
		if len(in.Images) > 0 {
			images := make([]*image.Image, 0, len(in.Images))
			for i, img := range in.Images {
				images = append(images, &image.Image{
					TargetType:  image.TargetTypeCatApply,
					TargetID:    apply.ID,
					URL:         img.Url,
					Sort:        img.Sort,
					IsCover:     img.IsCover,
					Description: img.Description,
					UploaderID:  in.ApplicantUserId,
				})

				// 如果请求里没指定 cover，默认第一张是封面
				if i == 0 && !hasCoverImage(in.Images) {
					images[i].IsCover = true
				}
			}

			if err := l.svcCtx.Repo.Image.BatchCreate(l.ctx, images, tx); err != nil {
				logger.Error("创建图片记录失败", zap.Error(err))
				return errorx.WrapDBInsert("创建图片记录失败", err)
			}
		}

		// ========== 3. 创建标签关联 ==========
		if len(in.TagIds) > 0 {
			relations := make([]*tag.TagRelation, 0, len(in.TagIds))
			for i, tagID := range in.TagIds {
				relations = append(relations, &tag.TagRelation{
					TagID:      tagID,
					TargetID:   apply.ID,
					TargetType: tag.TargetTypeCatApply,
					Sort:       int32(i),
					CreatedBy:  in.ApplicantUserId,
				})
			}

			if err := l.svcCtx.Repo.Tag.BatchAddTagRelations(l.ctx, relations, tx); err != nil {
				logger.Error("创建标签关联失败", zap.Error(err))
				return errorx.WrapDBInsert("创建标签关联失败", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logger.Info("志愿者申请创建猫咪档案成功", zap.Uint64("applyID", applyID))

	return &v1.ApplyCreateCatResponse{
		ApplyId: applyID,
	}, nil
}

// hasCoverImage 检查是否已有指定封面
func hasCoverImage(images []*v1.ImageItem) bool {
	for _, img := range images {
		if img.IsCover {
			return true
		}
	}
	return false
}
