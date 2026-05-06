package logic

import (
	"context"
	"math"

	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
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
	clustered := make([]bool, n) // 是否已经归属某个cluster

	var hotspots []*hotspot

	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}

		visited[i] = true
		neighbors := l.findNeighbors(cats, i, radiusKm)

		// 不是核心点
		if len(neighbors) < minPts {
			continue
		}

		// 创建 cluster
		clusterPoints := l.expandCluster(cats, i, neighbors, visited, clustered, radiusKm, minPts)

		// 生成 hotspot
		hs := l.buildHotspot(cats, clusterPoints)
		hotspots = append(hotspots, hs)
	}

	return hotspots
}

func (l *DetectHotspotsLogic) expandCluster(
	cats []*cat.Cat,
	pointIdx int,
	neighbors []int,
	visited []bool,
	clustered []bool,
	radiusKm float64,
	minPts int,
) []int {

	cluster := make([]int, 0)
	queue := append([]int{}, neighbors...)

	// 当前点加入 cluster
	cluster = append(cluster, pointIdx)
	clustered[pointIdx] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if !visited[curr] {
			visited[curr] = true

			currNeighbors := l.findNeighbors(cats, curr, radiusKm)
			if len(currNeighbors) >= minPts {
				queue = append(queue, currNeighbors...)
			}
		}

		// 加入 cluster（避免重复）
		if !clustered[curr] {
			clustered[curr] = true
			cluster = append(cluster, curr)
		}
	}

	return cluster
}

func (l *DetectHotspotsLogic) buildHotspot(cats []*cat.Cat, indices []int) *hotspot {
	var sumLat, sumLng float64
	catIDs := make([]uint64, 0, len(indices))

	for _, idx := range indices {
		sumLat += cats[idx].Latitude
		sumLng += cats[idx].Longitude
		catIDs = append(catIDs, cats[idx].ID)
	}

	return &hotspot{
		centerLat: sumLat / float64(len(indices)),
		centerLng: sumLng / float64(len(indices)),
		catIDs:    catIDs,
		density:   len(indices),
	}
}

func (l *DetectHotspotsLogic) findNeighbors(cats []*cat.Cat, centerIdx int, radiusKm float64) []int {
	neighbors := []int{}
	center := cats[centerIdx]

	for i := 0; i < len(cats); i++ {
		dist := l.haversine(
			center.Latitude, center.Longitude,
			cats[i].Latitude, cats[i].Longitude,
		)
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
