package logic

import (
	"context"
	"math"

	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/zeromicro/go-zero/core/logx"
)

type DetectHotspotsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDetectHotspotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetectHotspotsLogic {
	return &DetectHotspotsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DetectHotspotsLogic) DetectHotspots(in *v1.DetectHotspotsRequest) (*v1.DetectHotspotsResponse, error) {
	// 1. 获取所有有坐标的猫咪
	cats, err := l.svcCtx.Repo.Cat.GetAllCatsWithLocation(l.ctx)
	if err != nil {
		return nil, err
	}

	if len(cats) == 0 {
		return &v1.DetectHotspotsResponse{
			Hotspots:      []*v1.HotspotInfo{},
			TotalHotspots: 0,
		}, nil
	}

	// 2. DBSCAN 聚类
	hotspots := l.dbscanCluster(cats, in.RadiusKm, int(in.MinCats))

	// 3. 转换响应
	resp := &v1.DetectHotspotsResponse{
		Hotspots:      make([]*v1.HotspotInfo, len(hotspots)),
		TotalHotspots: int32(len(hotspots)),
	}

	for i, hs := range hotspots {
		resp.Hotspots[i] = &v1.HotspotInfo{
			CenterLat: hs.centerLat,
			CenterLng: hs.centerLng,
			CatIds:    hs.catIDs,
			Density:   int32(hs.density),
			RadiusKm:  in.RadiusKm,
		}
	}

	return resp, nil
}

type hotspot struct {
	centerLat, centerLng float64
	catIDs               []uint64
	density              int
}

func (l *DetectHotspotsLogic) dbscanCluster(cats []*cat.Cat, radiusKm float64, minPts int) []*hotspot {
	n := len(cats)
	visited := make([]bool, n)
	hotspots := make([]*hotspot, 0)

	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}

		// 查找邻居
		neighbors := l.findNeighbors(cats, i, radiusKm)

		if len(neighbors) >= minPts {
			// 形成核心点，构建聚类
			cluster := &hotspot{
				catIDs:  make([]uint64, 0),
				density: len(neighbors),
			}

			var sumLat, sumLng float64
			for _, idx := range neighbors {
				visited[idx] = true
				cluster.catIDs = append(cluster.catIDs, cats[idx].ID)
				sumLat += cats[idx].Latitude
				sumLng += cats[idx].Longitude
			}

			cluster.centerLat = sumLat / float64(len(neighbors))
			cluster.centerLng = sumLng / float64(len(neighbors))
			hotspots = append(hotspots, cluster)
		}
	}

	return hotspots
}

func (l *DetectHotspotsLogic) findNeighbors(cats []*cat.Cat, centerIdx int, radiusKm float64) []int {
	neighbors := []int{centerIdx}
	center := cats[centerIdx]

	for i := 0; i < len(cats); i++ {
		if i == centerIdx {
			continue
		}

		dist := l.haversine(center.Latitude, center.Longitude,
			cats[i].Latitude, cats[i].Longitude)
		if dist <= radiusKm {
			neighbors = append(neighbors, i)
		}
	}

	return neighbors
}

func (l *DetectHotspotsLogic) haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
