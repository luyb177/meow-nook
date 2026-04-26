// internal/logic/direct_create_cat_logic.go
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

type DirectCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDirectCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DirectCreateCatLogic {
	return &DirectCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DirectCreateCatLogic) DirectCreateCat(in *v1.DirectCreateCatRequest) (*v1.DirectCreateCatResponse, error) {
	var catID uint64
	var catCode string

	err := l.svcCtx.Repo.WithTx(l.ctx, func(tx *gorm.DB) error {
		// ========== 1. 创建正式猫咪档案 ==========
		catCode = in.CatCode
		if catCode == "" {
			var err error
			catCode, err = l.svcCtx.Repo.Sequence.GenerateCatCode(l.ctx)
			if err != nil {
				return errorx.WrapRedisSet("生成猫咪编号失败", err)
			}
		}

		newCat := &cat.Cat{
			CatCode:                 catCode,
			Name:                    in.Name,
			Breed:                   in.Breed,
			Color:                   in.Color,
			Gender:                  in.Gender,
			BodySize:                in.BodySize,
			AgeStage:                in.AgeStage,
			Description:             in.Description,
			DiscoveryAddress:        in.DiscoveryAddress,
			Longitude:               in.Longitude,
			Latitude:                in.Latitude,
			IsVaccinated:            in.IsVaccinated,
			IsHealthy:               in.IsHealthy,
			NeedMedicalIntervention: in.NeedMedicalIntervention,
			SterilizationStatus:     in.SterilizationStatus,
			AdoptionStatus:          cat.CatAdoptionStatusPending,
			ApplyID:                 0, // 管理员直接创建，无来源申请
			CreatorID:               in.OperatorId,
		}

		if err := l.svcCtx.Repo.Cat.CreateCat(l.ctx, newCat, tx); err != nil {
			return errorx.WrapDBInsert("创建猫咪档案失败", err)
		}

		catID = newCat.ID

		// ========== 2. 创建图片 ==========
		if len(in.Images) > 0 {
			images := make([]*image.Image, 0, len(in.Images))
			for i, img := range in.Images {
				images = append(images, &image.Image{
					TargetType:  image.TargetTypeCatProfile,
					TargetID:    catID,
					URL:         img.Url,
					Sort:        img.Sort,
					IsCover:     img.IsCover,
					Description: img.Description,
					UploaderID:  in.OperatorId,
				})

				if i == 0 && !hasCoverImage(in.Images) {
					images[i].IsCover = true
				}
			}

			if err := l.svcCtx.Repo.Image.BatchCreate(l.ctx, images, tx); err != nil {
				return errorx.WrapDBInsert("创建图片记录失败", err)
			}

			// ========== 3. 设置头像 ==========
			cover, err := l.svcCtx.Repo.Image.GetCover(
				l.ctx,
				image.TargetTypeCatProfile,
				catID,
				tx,
			)
			if err == nil && cover != nil {
				if err := l.svcCtx.Repo.Cat.UpdateCat(
					l.ctx,
					catID,
					map[string]any{"avatar": cover.URL},
					tx,
				); err != nil {
					return errorx.WrapDBUpdate("设置猫咪头像失败", err)
				}
			}
		}

		// ========== 4. 创建标签关联 ==========
		if len(in.TagIds) > 0 {
			relations := make([]*tag.TagRelation, 0, len(in.TagIds))
			for i, tagID := range in.TagIds {
				relations = append(relations, &tag.TagRelation{
					TagID:      tagID,
					TargetID:   catID,
					TargetType: tag.TargetTypeCatProfile,
					Sort:       int32(i),
					CreatedBy:  in.OperatorId,
				})
			}

			if err := l.svcCtx.Repo.Tag.BatchAddTagRelations(
				l.ctx,
				relations,
				tx,
			); err != nil {
				return errorx.WrapDBInsert("创建标签关联失败", err)
			}
		}

		return nil
	})

	if err != nil {
		logger.Error("管理员直接创建猫咪档案失败", zap.Error(err))
		return nil, err
	}

	logger.Info("管理员直接创建猫咪档案成功", zap.Uint64("catID", catID), zap.String("catCode", catCode))

	return &v1.DirectCreateCatResponse{
		CatId:   catID,
		CatCode: catCode,
	}, nil
}
