package logic

import (
	"context"
	"errors"

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

type ApproveCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveCreateCatLogic {
	return &ApproveCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveCreateCatLogic) ApproveCreateCat(in *v1.ApproveCreateCatRequest) (*v1.ApproveCreateCatResponse, error) {
	var catID uint64
	var catCode string

	err := l.svcCtx.Repo.WithTx(l.ctx, func(tx *gorm.DB) error {
		// ========== 1. 查询并锁定申请（防并发） ==========
		apply, err := l.svcCtx.Repo.Cat.GetApplyByIDForUpdate(l.ctx, in.ApplyId, tx)
		if err != nil {
			return errorx.WrapDBQuery("查询申请失败", err)
		}

		if apply.Status != cat.ApplyStatusPending {
			return errorx.Wrap(errorx.CodeBadRequest, "申请状态异常，无法审核", errors.New("invalid application status"))
		}

		// ========== 2. 创建正式猫咪档案 ==========
		catCode = in.CatCode
		if catCode == "" {
			// 可以自动生成编号，例如：CAT-20240101-001
			catCode, err = l.svcCtx.Repo.Sequence.GenerateCatCode(l.ctx)
			if err != nil {
				return errorx.Wrap(errorx.CodeRedisSetFailed, "生成猫咪编号失败", err)
			}
		}

		newCat := &cat.Cat{
			CatCode:                 catCode,
			Name:                    apply.Name,
			Gender:                  apply.Gender,
			BodySize:                apply.BodySize,
			AgeStage:                apply.AgeStage,
			Description:             apply.Description,
			DiscoveryAddress:        apply.DiscoveryAddress,
			Longitude:               apply.Longitude,
			Latitude:                apply.Latitude,
			Breed:                   in.Breed,
			Color:                   in.Color,
			IsVaccinated:            in.IsVaccinated,
			IsHealthy:               in.IsHealthy,
			NeedMedicalIntervention: in.NeedMedicalIntervention,
			SterilizationStatus:     in.SterilizationStatus,
			AdoptionStatus:          "pending",
			ApplyID:                 apply.ID,
			CreatorID:               in.OperatorId,
		}

		if err := l.svcCtx.Repo.Cat.CreateCat(l.ctx, newCat, tx); err != nil {
			return errorx.WrapDBInsert("创建猫咪档案失败", err)
		}

		catID = newCat.ID

		// ========== 3. 更新申请状态 ==========
		if err := l.svcCtx.Repo.Cat.ApproveApply(
			l.ctx,
			apply.ID,
			in.OperatorId,
			catID,
			tx,
		); err != nil {
			return errorx.WrapDBUpdate("更新申请状态失败", err)
		}

		// ========== 4. 复制图片到正式档案 ==========
		if err := l.svcCtx.Repo.Image.CopyByTarget(
			l.ctx,
			image.TargetTypeCatApply, apply.ID,
			image.TargetTypeCatProfile, catID,
			apply.ApplicantUserID,
			tx,
		); err != nil {
			return errorx.WrapDBInsert("复制图片失败", err)
		}

		// ========== 5. 设置头像 ==========
		cover, err := l.svcCtx.Repo.Image.GetCover(
			l.ctx,
			image.TargetTypeCatProfile,
			catID,
			tx,
		)
		if err != nil {
			logger.Error("获取封面图失败", zap.Error(err))
		}
		if cover != nil {
			if err := l.svcCtx.Repo.Cat.UpdateCat(
				l.ctx,
				catID,
				map[string]any{"avatar": cover.URL},
				tx,
			); err != nil {
				return errorx.WrapDBUpdate("设置头像失败", err)
			}
		}

		// ========== 6. 迁移标签到正式档案 ==========
		// 6a. 查询申请时的标签
		applyTags, err := l.svcCtx.Repo.Tag.GetTagsByTarget(
			l.ctx,
			tag.TargetTypeCatApply,
			apply.ID,
			tx,
		)
		if err != nil {
			logger.Error("查询申请标签失败", zap.Error(err))
		}

		// 6b. 合并：申请标签 + 审核时追加的标签
		allTagIDs := make([]uint64, 0)
		for _, t := range applyTags {
			allTagIDs = append(allTagIDs, t.ID)
		}
		for _, extraID := range in.ExtraTagIds {
			if !contains(allTagIDs, extraID) {
				allTagIDs = append(allTagIDs, extraID)
			}
		}

		// 6c. 写入正式档案标签
		if len(allTagIDs) > 0 {
			relations := make([]*tag.TagRelation, 0, len(allTagIDs))
			for i, tagID := range allTagIDs {
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
		logger.Error("审核通过失败", zap.Uint64("apply_id", in.ApplyId), zap.Error(err))
		return nil, err
	}

	logger.Info("审核通过成功", zap.Uint64("apply_id", in.ApplyId), zap.Uint64("cat_id", catID))

	return &v1.ApproveCreateCatResponse{
		CatId:   catID,
		CatCode: catCode,
		Status:  cat.ApplyStatusApproved,
	}, nil
}

func contains(slice []uint64, target uint64) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}
