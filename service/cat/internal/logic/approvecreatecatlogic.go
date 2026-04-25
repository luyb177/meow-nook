package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApproveCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveCreateCatLogic {
	return &ApproveCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApproveCreateCatLogic) ApproveCreateCat(in *v1.ApproveCreateCatRequest) (*v1.ApproveCreateCatResponse, error) {
	// TODO(permission): check that in.OperatorId has admin role

	apply, err := l.svcCtx.Repo.Cat.GetCatCreateApplyByID(l.ctx, in.ApplyId)
	if err != nil {
		return nil, errorx.ErrCatNotFound
	}
	if apply.Status != "pending" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "申请状态不是待审核", nil)
	}

	newCat := &cat.Cat{
		CatCode:          fmt.Sprintf("CAT-%06d", apply.ID),
		Name:             apply.Name,
		Gender:           apply.Gender,
		BodySize:         apply.BodySize,
		AgeStage:         apply.AgeStage,
		Description:      apply.Description,
		DiscoveryAddress: apply.DiscoveryAddress,
		Longitude:        apply.Longitude,
		Latitude:         apply.Latitude,
		AdoptionStatus:   "pending",
	}
	if err := l.svcCtx.Repo.Cat.CreateCat(l.ctx, newCat); err != nil {
		logger.Error("ApproveCreateCat: create cat failed")
		return nil, errorx.WrapInternal("创建猫咪档案失败", err)
	}

	// Store images linked to the new cat
	var imageURLs []string
	if apply.ImageURLs != "" {
		_ = json.Unmarshal([]byte(apply.ImageURLs), &imageURLs)
	}
	if len(imageURLs) > 0 {
		images := make([]*cat.CatImage, 0, len(imageURLs))
		for i, u := range imageURLs {
			images = append(images, &cat.CatImage{
				CatID:      newCat.ID,
				ImageURL:   u,
				ImageType:  "normal",
				Sort:       int32(i),
				UploaderID: apply.ApplicantUserID,
			})
		}
		if err := l.svcCtx.Repo.Cat.CreateCatImages(l.ctx, images); err != nil {
			logger.Error("ApproveCreateCat: create images failed")
			// non-fatal: cat record already created
		}
	}

	if err := l.svcCtx.Repo.Cat.UpdateCatCreateApplyStatus(l.ctx, apply.ID, "approved", "", in.OperatorId); err != nil {
		logger.Error("ApproveCreateCat: update apply status failed")
		return nil, errorx.WrapInternal("更新申请状态失败", err)
	}

	return &v1.ApproveCreateCatResponse{CatId: newCat.ID}, nil
}
