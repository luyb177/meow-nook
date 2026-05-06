package algorithm

import (
	"context"

	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type DetectHotspotsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetectHotspotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetectHotspotsLogic {
	return &DetectHotspotsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetectHotspotsLogic) DetectHotspots(req *types.DetectHotspotsReq) (*types.DetectHotspotsResp, error) {
	// 设置默认值
	radiusKm := req.RadiusKm
	if radiusKm <= 0 {
		radiusKm = 0.5 // 默认500米
	}

	minCats := req.MinCats
	if minCats <= 0 {
		minCats = 3
	}

	// 调用 cat 服务
	resp, err := l.svcCtx.CatRPC.DetectHotspots(l.ctx, &catpb.DetectHotspotsRequest{
		RadiusKm: radiusKm,
		MinCats:  int32(minCats),
	})
	if err != nil {
		return nil, err
	}

	// 转换响应
	hotspots := make([]types.HotspotInfo, len(resp.Hotspots))
	for i, hs := range resp.Hotspots {
		hotspots[i] = types.HotspotInfo{
			CenterLat: hs.CenterLat,
			CenterLng: hs.CenterLng,
			CatIDs:    hs.CatIds,
			Density:   int(hs.Density),
			RadiusKm:  hs.RadiusKm,
		}
	}

	return &types.DetectHotspotsResp{
		Hotspots:      hotspots,
		TotalHotspots: int(resp.TotalHotspots),
	}, nil
}
