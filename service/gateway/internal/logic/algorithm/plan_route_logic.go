package algorithm

import (
	"context"

	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type PlanRouteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPlanRouteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlanRouteLogic {
	return &PlanRouteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PlanRouteLogic) PlanRoute(req *types.PlanRouteReq) (*types.PlanRouteResp, error) {
	// 调用 cat 服务
	resp, err := l.svcCtx.CatRPC.PlanRoute(l.ctx, &catpb.PlanRouteRequest{
		VolunteerLat: req.VolunteerLat,
		VolunteerLng: req.VolunteerLng,
		CatIds:       req.CatIDs,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应
	waypoints := make([]types.Waypoint, len(resp.Waypoints))
	for i, wp := range resp.Waypoints {
		waypoints[i] = types.Waypoint{
			ID:    wp.Id,
			Type:  wp.Type,
			Lat:   wp.Lat,
			Lng:   wp.Lng,
			Order: int(wp.Order),
		}
	}

	return &types.PlanRouteResp{
		Waypoints:       waypoints,
		TotalDistanceKm: resp.TotalDistanceKm,
		TotalStops:      int(resp.TotalStops),
	}, nil
}
