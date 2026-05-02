package logic

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"go.uber.org/zap"
)

type ListCatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCatsLogic {
	return &ListCatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCatsLogic) ListCats(in *v1.ListCatsReq) (*v1.ListCatsResp, error) {
	// 参数校验和默认值设置
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}

	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 构建筛选条件
	filter := cat.CatListFilter{
		Keyword:             in.Keyword,
		Name:                in.Name,
		CatCode:             in.CatCode,
		Breed:               in.Breed,
		Color:               in.Color,
		Gender:              in.Gender,
		BodySize:            in.BodySize,
		AgeStage:            in.AgeStage,
		SterilizationStatus: in.SterilizationStatus,
		AdoptionStatus:      in.AdoptionStatus,
		AdopterID:           in.AdopterId,
		CreatorID:           in.CreatorId,
		ApplyID:             in.ApplyId,
		Page:                page,
		PageSize:            pageSize,
		SortBy:              in.SortBy,
		SortOrder:           in.SortOrder,
	}

	// Bool 类型处理
	if in.IsVaccinated != nil {
		filter.IsVaccinated = in.IsVaccinated
	}
	if in.IsHealthy != nil {
		filter.IsHealthy = in.IsHealthy
	}
	if in.NeedMedicalIntervention != nil {
		filter.NeedMedicalIntervention = in.NeedMedicalIntervention
	}

	// 时间字段处理
	if in.FoundAtStart != "" {
		if t, err := time.Parse(time.RFC3339, in.FoundAtStart); err == nil {
			filter.FoundAtStart = &t
		}
	}
	if in.FoundAtEnd != "" {
		if t, err := time.Parse(time.RFC3339, in.FoundAtEnd); err == nil {
			filter.FoundAtEnd = &t
		}
	}
	if in.CreatedAtStart != "" {
		if t, err := time.Parse(time.RFC3339, in.CreatedAtStart); err == nil {
			filter.CreatedAtStart = &t
		}
	}
	if in.CreatedAtEnd != "" {
		if t, err := time.Parse(time.RFC3339, in.CreatedAtEnd); err == nil {
			filter.CreatedAtEnd = &t
		}
	}

	var cats []*cat.Cat
	var distances []float64
	var total int64
	var err error

	// 判断是否使用附近查询
	if in.Nearby != nil && in.Nearby.Latitude != 0 && in.Nearby.Longitude != 0 {
		radius := in.Nearby.Radius
		if radius <= 0 {
			radius = 5 // 默认5公里
		}
		if radius > 50 {
			radius = 50 // 最大50公里
		}

		nearFilter := &cat.NearFilter{
			Latitude:  in.Nearby.Latitude,
			Longitude: in.Nearby.Longitude,
			Radius:    radius,
		}

		cats, distances, total, err = l.svcCtx.Repo.Cat.ListCatsWithNearby(l.ctx, filter, nearFilter)
		if err != nil {
			logger.Error("ListCatsWithNearby error", zap.Error(err))
			return nil, errorx.WrapDBQuery("查询附近猫咪失败: %v", err)
		}
	} else {
		// 普通查询
		cats, total, err = l.svcCtx.Repo.Cat.ListCats(l.ctx, filter)
		if err != nil {
			logger.Error("ListCats error", zap.Error(err))
			return nil, errorx.WrapDBQuery("查询猫咪列表失败: %v", err)
		}
	}

	// 构建响应
	resp := &v1.ListCatsResp{
		Total: total,
		List:  make([]*v1.CatInfo, 0, len(cats)),
	}

	for i, c := range cats {
		catInfo := &v1.CatInfo{
			Id:                      c.ID,
			CatCode:                 c.CatCode,
			Name:                    c.Name,
			Breed:                   c.Breed,
			Color:                   c.Color,
			Gender:                  c.Gender,
			BodySize:                c.BodySize,
			AgeStage:                c.AgeStage,
			Weight:                  c.Weight,
			Character:               c.Character,
			Avatar:                  c.Avatar,
			Description:             c.Description,
			DiscoveryAddress:        c.DiscoveryAddress,
			Longitude:               c.Longitude,
			Latitude:                c.Latitude,
			IsVaccinated:            c.IsVaccinated,
			IsHealthy:               c.IsHealthy,
			NeedMedicalIntervention: c.NeedMedicalIntervention,
			SterilizationStatus:     c.SterilizationStatus,
			AdoptionStatus:          c.AdoptionStatus,
			AdopterId:               c.AdopterID,
			ApplyId:                 c.ApplyID,
			CreatorId:               c.CreatorID,
		}

		// 时间字段转换
		if c.BirthDate != nil {
			catInfo.BirthDate = c.BirthDate.Format(time.RFC3339)
		}
		if c.FoundAt != nil {
			catInfo.FoundAt = c.FoundAt.Format(time.RFC3339)
		}
		if c.LastMedicalCheckAt != nil {
			catInfo.LastMedicalCheckAt = c.LastMedicalCheckAt.Format(time.RFC3339)
		}
		if c.AdoptedAt != nil {
			catInfo.AdoptedAt = c.AdoptedAt.Format(time.RFC3339)
		}
		catInfo.CreatedAt = c.CreatedAt.Format(time.RFC3339)
		catInfo.UpdatedAt = c.UpdatedAt.Format(time.RFC3339)

		// 设置距离（如果有）
		if len(distances) > i {
			catInfo.Distance = distances[i]
		}

		resp.List = append(resp.List, catInfo)
	}

	return resp, nil
}
