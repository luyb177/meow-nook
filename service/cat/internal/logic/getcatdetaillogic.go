package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCatDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCatDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCatDetailLogic {
	return &GetCatDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 猫咪详情
func (l *GetCatDetailLogic) GetCatDetail(in *v1.GetCatDetailRequest) (*v1.GetCatDetailResponse, error) {
	c, err := l.svcCtx.Repo.Cat.GetCatByID(l.ctx, in.CatId)
	if err != nil {
		return nil, errorx.ErrCatNotFound
	}

	images, err := l.svcCtx.Repo.Cat.ListCatImages(l.ctx, c.ID)
	if err != nil {
		logger.Error("GetCatDetail: list images failed")
		return nil, errorx.WrapInternal("查询图片失败", err)
	}

	imageURLs := make([]string, 0, len(images))
	for _, img := range images {
		imageURLs = append(imageURLs, img.ImageURL)
	}

	return &v1.GetCatDetailResponse{
		CatId:                   c.ID,
		CatCode:                 c.CatCode,
		Name:                    c.Name,
		Gender:                  c.Gender,
		BodySize:                c.BodySize,
		AgeStage:                c.AgeStage,
		SterilizationStatus:     c.SterilizationStatus,
		AdoptionStatus:          c.AdoptionStatus,
		Avatar:                  c.Avatar,
		Description:             c.Description,
		DiscoveryAddress:        c.DiscoveryAddress,
		Longitude:               c.Longitude,
		Latitude:                c.Latitude,
		IsVaccinated:            c.IsVaccinated,
		IsHealthy:               c.IsHealthy,
		NeedMedicalIntervention: c.NeedMedicalIntervention,
		ImageUrls:               imageURLs,
		CreatedAt:               timestamppb.New(c.CreatedAt),
		UpdatedAt:               timestamppb.New(c.UpdatedAt),
	}, nil
}
