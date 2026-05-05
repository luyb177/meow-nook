package logic

import (
	"context"
	"math"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/zeromicro/go-zero/core/logx"
)

type PlanRouteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPlanRouteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlanRouteLogic {
	return &PlanRouteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PlanRouteLogic) PlanRoute(in *v1.PlanRouteRequest) (*v1.PlanRouteResponse, error) {
	// 1. 获取猫咪位置
	cats, err := l.svcCtx.Repo.Cat.GetCatsByIDs(l.ctx, in.CatIds)
	if err != nil {
		return nil, err
	}

	if len(cats) == 0 {
		return &v1.PlanRouteResponse{
			Waypoints:       []*v1.Waypoint{},
			TotalDistanceKm: 0,
			TotalStops:      0,
		}, nil
	}

	// 2. 构建位置点
	start := location{lat: in.VolunteerLat, lng: in.VolunteerLng, id: 0, typ: "volunteer"}
	targets := make([]location, len(cats))
	for i, c := range cats {
		targets[i] = location{lat: c.Latitude, lng: c.Longitude, id: c.ID, typ: "cat"}
	}

	// 3. 最近邻算法
	route := l.nearestNeighbor(start, targets)

	// 4. 转换响应
	waypoints := make([]*v1.Waypoint, len(route.waypoints))
	for i, wp := range route.waypoints {
		waypoints[i] = &v1.Waypoint{
			Id:    wp.id,
			Type:  wp.typ,
			Lat:   wp.lat,
			Lng:   wp.lng,
			Order: int32(i),
		}
	}

	return &v1.PlanRouteResponse{
		Waypoints:       waypoints,
		TotalDistanceKm: route.totalDistance,
		TotalStops:      int32(len(route.waypoints) - 1), // 减掉起点
	}, nil
}

type location struct {
	id       uint64
	lat, lng float64
	typ      string
}

type route struct {
	waypoints     []location
	totalDistance float64
}

func (l *PlanRouteLogic) nearestNeighbor(start location, targets []location) *route {
	if len(targets) == 0 {
		return &route{
			waypoints:     []location{start},
			totalDistance: 0,
		}
	}

	visited := make([]bool, len(targets))
	result := &route{
		waypoints: []location{start},
	}
	current := start

	for len(result.waypoints)-1 < len(targets) {
		nearestIdx := -1
		minDist := math.MaxFloat64

		for i, target := range targets {
			if visited[i] {
				continue
			}
			dist := l.haversine(current.lat, current.lng, target.lat, target.lng)
			if dist < minDist {
				minDist = dist
				nearestIdx = i
			}
		}

		if nearestIdx == -1 {
			break
		}

		result.waypoints = append(result.waypoints, targets[nearestIdx])
		result.totalDistance += minDist
		current = targets[nearestIdx]
		visited[nearestIdx] = true
	}

	return result
}

func (l *PlanRouteLogic) haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
