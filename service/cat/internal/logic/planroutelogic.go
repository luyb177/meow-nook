package logic

import (
	"context"
	"math"

	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"go.uber.org/zap"
)

type PlanRouteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPlanRouteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlanRouteLogic {
	return &PlanRouteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PlanRouteLogic) PlanRoute(in *v1.PlanRouteRequest) (*v1.PlanRouteResponse, error) {
	cats, err := l.svcCtx.Repo.Cat.GetCatsByIDs(l.ctx, in.CatIds)
	if err != nil {
		return nil, err
	}

	logger.Info("规划路线请求", zap.Float64("volunteer_lat", in.VolunteerLat), zap.Float64("volunteer_lng", in.VolunteerLng), zap.Int("cat_count", len(cats)))
	// 过滤无效坐标
	validCats := make([]*cat.Cat, 0)
	for _, c := range cats {
		if c.Latitude != 0 && c.Longitude != 0 {
			validCats = append(validCats, c)
		}
	}

	if len(validCats) == 0 {
		return &v1.PlanRouteResponse{}, nil
	}

	start := location{
		id:  0,
		lat: in.VolunteerLat,
		lng: in.VolunteerLng,
		typ: "volunteer",
	}

	targets := make([]location, 0, len(validCats))
	for _, c := range validCats {
		targets = append(targets, location{
			id:  c.ID,
			lat: c.Latitude,
			lng: c.Longitude,
			typ: "cat",
		})
	}

	// 最近邻生成初始路径
	path := l.nearestNeighborPath(start, targets)

	// 2-opt 优化路径
	path = l.twoOpt(path)

	// 计算总距离
	totalDist := l.calcTotalDistance(path)

	// 转换响应
	waypoints := make([]*v1.Waypoint, len(path))
	for i, p := range path {
		waypoints[i] = &v1.Waypoint{
			Id:    p.id,
			Type:  p.typ,
			Lat:   p.lat,
			Lng:   p.lng,
			Order: int32(i),
		}
	}

	return &v1.PlanRouteResponse{
		Waypoints:       waypoints,
		TotalDistanceKm: totalDist,
		TotalStops:      int32(len(path) - 1),
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

func (l *PlanRouteLogic) nearestNeighborPath(start location, targets []location) []location {
	path := []location{start}
	visited := make([]bool, len(targets))

	current := start

	for len(path)-1 < len(targets) {
		minDist := math.MaxFloat64
		idx := -1

		for i, t := range targets {
			if visited[i] {
				continue
			}
			dist := l.haversine(current.lat, current.lng, t.lat, t.lng)
			if dist < minDist {
				minDist = dist
				idx = i
			}
		}

		if idx == -1 {
			break
		}

		path = append(path, targets[idx])
		current = targets[idx]
		visited[idx] = true
	}

	return path
}

func (l *PlanRouteLogic) twoOpt(path []location) []location {
	improved := true
	n := len(path)

	for improved {
		improved = false

		for i := 1; i < n-2; i++ {
			for j := i + 1; j < n-1; j++ {

				a, b := path[i-1], path[i]
				c, d := path[j], path[j+1]

				oldDist := l.haversine(a.lat, a.lng, b.lat, b.lng) +
					l.haversine(c.lat, c.lng, d.lat, d.lng)

				newDist := l.haversine(a.lat, a.lng, c.lat, c.lng) +
					l.haversine(b.lat, b.lng, d.lat, d.lng)

				if newDist < oldDist {
					// 反转路径
					l.reverse(path, i, j)
					improved = true
				}
			}
		}
	}

	return path
}

func (l *PlanRouteLogic) reverse(path []location, i, j int) {
	for i < j {
		path[i], path[j] = path[j], path[i]
		i++
		j--
	}
}

func (l *PlanRouteLogic) calcTotalDistance(path []location) float64 {
	total := 0.0
	for i := 1; i < len(path); i++ {
		total += l.haversine(
			path[i-1].lat, path[i-1].lng,
			path[i].lat, path[i].lng,
		)
	}
	return total
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
