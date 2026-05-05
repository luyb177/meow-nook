package logic

import (
	"context"
	"math"
	"sort"

	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPriorityCatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPriorityCatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPriorityCatsLogic {
	return &GetPriorityCatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPriorityCatsLogic) GetPriorityCats(in *v1.GetPriorityCatsRequest) (*v1.GetPriorityCatsResponse, error) {
	// 1. 获取所有待领养或需要救助的猫咪
	filter := cat.CatListFilter{
		AdoptionStatus: "pending", // 待领养状态
		Page:           1,
		PageSize:       1000, // 获取足够多的猫咪
	}

	cats, _, err := l.svcCtx.Repo.Cat.ListCats(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(cats) == 0 {
		return &v1.GetPriorityCatsResponse{
			Cats:  []*v1.PriorityCatInfo{},
			Total: 0,
		}, nil
	}

	// 2. 获取猫咪ID列表
	catIDs := make([]uint64, len(cats))
	for i, c := range cats {
		catIDs[i] = c.ID
	}

	// 3. 获取报告次数（7天内）
	reportCounts, err := l.svcCtx.Repo.Cat.GetCatReportCounts(l.ctx, catIDs, 7)
	if err != nil {
		reportCounts = make(map[uint64]int) // 降级处理
	}

	// 4. 计算优先级分数
	scores := make([]*priorityScore, 0, len(cats))
	for _, c := range cats {
		score := l.calculatePriority(c, reportCounts[c.ID])
		scores = append(scores, score)
	}

	// 5. 排序
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// 6. 取前N个
	topN := int(in.TopN)
	if topN <= 0 || topN > len(scores) {
		topN = len(scores)
	}

	// 7. 构建响应
	resp := &v1.GetPriorityCatsResponse{
		Cats:  make([]*v1.PriorityCatInfo, topN),
		Total: int32(topN),
	}

	for i := 0; i < topN; i++ {
		s := scores[i]
		resp.Cats[i] = &v1.PriorityCatInfo{
			CatId:     s.catID,
			Score:     s.score,
			Details:   s.details,
			CatName:   s.catName,
			CatCode:   s.catCode,
			Longitude: s.longitude,
			Latitude:  s.latitude,
		}
	}

	return resp, nil
}

type priorityScore struct {
	catID     uint64
	score     float64
	details   map[string]float64
	catName   string
	catCode   string
	longitude float64
	latitude  float64
}

func (l *GetPriorityCatsLogic) calculatePriority(c *cat.Cat, reportCount int) *priorityScore {
	weights := map[string]float64{
		"health_risk": 0.35,
		"age_factor":  0.25,
		"environment": 0.20,
		"report_freq": 0.20,
	}

	details := make(map[string]float64)

	// 1. 健康风险 (0-100)
	healthScore := 0.0
	if c.NeedMedicalIntervention {
		healthScore += 50
	}
	if !c.IsHealthy {
		healthScore += 30
	}
	if !c.IsVaccinated {
		healthScore += 20
	}
	details["health_risk"] = healthScore

	// 2. 年龄因子 (0-100)
	ageScore := 0.0
	switch c.AgeStage {
	case "kitten":
		ageScore = 100
	case "young":
		ageScore = 70
	case "adult":
		ageScore = 40
	case "old":
		ageScore = 80
	default:
		ageScore = 50
	}
	details["age_factor"] = ageScore

	// 3. 环境危险度 (0-100)
	envScore := 0.0
	if c.NeedMedicalIntervention {
		envScore += 40
	}
	// 可以扩展：检查是否在危险区域（马路附近等）
	details["environment"] = envScore

	// 4. 报告频率 (0-100)
	reportScore := math.Min(float64(reportCount)*10, 100)
	details["report_freq"] = reportScore

	// 加权总分
	totalScore := weights["health_risk"]*healthScore +
		weights["age_factor"]*ageScore +
		weights["environment"]*envScore +
		weights["report_freq"]*reportScore

	return &priorityScore{
		catID:     c.ID,
		score:     totalScore,
		details:   details,
		catName:   c.Name,
		catCode:   c.CatCode,
		longitude: c.Longitude,
		latitude:  c.Latitude,
	}
}
