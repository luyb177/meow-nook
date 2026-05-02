package cat

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
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

// ListCats 获取小猫列表
func (l *ListCatsLogic) ListCats(req *types.ListCatsReq) (resp *types.ListCatsResp, err error) {
	// 记录请求日志
	logger.FromContext(l.ctx).Info("ListCats request",
		zap.String("keyword", req.Keyword),
		zap.Int32("page", req.Page),
		zap.Int32("page_size", req.PageSize),
	)

	// 参数校验
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 构建 RPC 请求
	rpcReq := &catpb.ListCatsReq{
		Keyword:             req.Keyword,
		Name:                req.Name,
		CatCode:             req.CatCode,
		Breed:               req.Breed,
		Color:               req.Color,
		Gender:              req.Gender,
		BodySize:            req.BodySize,
		AgeStage:            req.AgeStage,
		SterilizationStatus: req.SterilizationStatus,
		AdoptionStatus:      req.AdoptionStatus,
		AdopterId:           req.AdopterID,
		CreatorId:           req.CreatorID,
		ApplyId:             req.ApplyID,
		Page:                req.Page,
		PageSize:            req.PageSize,
		SortBy:              req.SortBy,
		SortOrder:           req.SortOrder,
		FoundAtStart:        req.FoundAtStart,
		FoundAtEnd:          req.FoundAtEnd,
		CreatedAtStart:      req.CreatedAtStart,
		CreatedAtEnd:        req.CreatedAtEnd,
	}

	// 处理 bool 类型
	if req.IsVaccinated != nil {
		rpcReq.IsVaccinated = req.IsVaccinated
	}
	if req.IsHealthy != nil {
		rpcReq.IsHealthy = req.IsHealthy
	}
	if req.NeedMedicalIntervention != nil {
		rpcReq.NeedMedicalIntervention = req.NeedMedicalIntervention
	}

	// 处理附近筛选
	if req.Nearby != nil {
		rpcReq.Nearby = &catpb.NearbyFilter{
			Latitude:  req.Nearby.Latitude,
			Longitude: req.Nearby.Longitude,
			Radius:    req.Nearby.Radius,
		}
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.CatRPC.ListCats(l.ctx, rpcReq)
	if err != nil {
		logger.FromContext(l.ctx).Error("ListCats RPC call failed",
			zap.Error(err),
		)
		return nil, errorx.FromGRPC(err)
	}

	// 转换响应数据
	resp = &types.ListCatsResp{
		Total: rpcResp.Total,
		List:  make([]types.CatInfo, 0, len(rpcResp.List)),
	}

	for _, cat := range rpcResp.List {
		catInfo := types.CatInfo{
			Id:                      cat.Id,
			CatCode:                 cat.CatCode,
			Name:                    cat.Name,
			Breed:                   cat.Breed,
			Color:                   cat.Color,
			Gender:                  cat.Gender,
			BodySize:                cat.BodySize,
			AgeStage:                cat.AgeStage,
			Weight:                  cat.Weight,
			Character:               cat.Character,
			Avatar:                  cat.Avatar,
			Description:             cat.Description,
			DiscoveryAddress:        cat.DiscoveryAddress,
			Longitude:               cat.Longitude,
			Latitude:                cat.Latitude,
			IsVaccinated:            cat.IsVaccinated,
			IsHealthy:               cat.IsHealthy,
			NeedMedicalIntervention: cat.NeedMedicalIntervention,
			SterilizationStatus:     cat.SterilizationStatus,
			AdoptionStatus:          cat.AdoptionStatus,
			AdopterId:               cat.AdopterId,
			ApplyId:                 cat.ApplyId,
			CreatorId:               cat.CreatorId,
			Distance:                cat.Distance,
		}

		// 转换时间字段（如果非空）
		if cat.BirthDate != "" {
			catInfo.BirthDate = cat.BirthDate
		}
		if cat.FoundAt != "" {
			catInfo.FoundAt = cat.FoundAt
		}
		if cat.LastMedicalCheckAt != "" {
			catInfo.LastMedicalCheckAt = cat.LastMedicalCheckAt
		}
		if cat.AdoptedAt != "" {
			catInfo.AdoptedAt = cat.AdoptedAt
		}
		if cat.CreatedAt != "" {
			catInfo.CreatedAt = cat.CreatedAt
		}
		if cat.UpdatedAt != "" {
			catInfo.UpdatedAt = cat.UpdatedAt
		}

		resp.List = append(resp.List, catInfo)
	}

	logger.FromContext(l.ctx).Info("ListCats response",
		zap.Int64("total", resp.Total),
		zap.Int("list_count", len(resp.List)),
	)

	return resp, nil
}
