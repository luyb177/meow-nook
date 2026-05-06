package cat

import (
	"context"
	"errors"
	"math"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

const (
	CatGenderUnknown = "unknown"
	CatGenderMale    = "male"
	CatGenderFemale  = "female"
)

const (
	CatBodySizeSmall  = "small"
	CatBodySizeMedium = "medium"
	CatBodySizeLarge  = "large"
)

const (
	CatAgeStageKitten = "kitten"
	CatAgeStageYoung  = "young"
	CatAgeStageAdult  = "adult"
	CatAgeStageOld    = "old"
)

const (
	CatSterilizationStatusSterilized   = "sterilized"
	CatSterilizationStatusUnsterilized = "unsterilized"
)

const (
	CatAdoptionStatusPending     = "pending"
	CatAdoptionStatusAdopted     = "adopted"
	CatAdoptionStatusUnavailable = "unavailable"
)

var (
	ErrCatCodeRequired        = errors.New("cat_code is required")
	ErrCatNameRequired        = errors.New("cat name is required")
	ErrInvalidCatID           = errors.New("invalid cat id")
	ErrInvalidAdoptionStatus  = errors.New("invalid adoption status")
	ErrAdopterIDRequired      = errors.New("adopter_id is required when adoption status is adopted")
	ErrCannotDeleteAdoptedCat = errors.New("cannot delete adopted cat")
)

type CatListFilter struct {
	Keyword string // name/cat_code/breed/color 模糊查询

	CatCode string
	Name    string
	Breed   string
	Color   string
	Gender  string

	BodySize string
	AgeStage string

	SterilizationStatus string
	AdoptionStatus      string

	AdopterID uint64
	CreatorID uint64
	ApplyID   uint64

	// bool 类型建议用指针，否则无法区分 false 和未传
	IsVaccinated            *bool
	IsHealthy               *bool
	NeedMedicalIntervention *bool

	FoundAtStart   *time.Time
	FoundAtEnd     *time.Time
	CreatedAtStart *time.Time
	CreatedAtEnd   *time.Time

	Page     int
	PageSize int

	// 排序字段，注意实现里要做白名单，避免 SQL 注入
	SortBy    string // id/created_at/found_at/updated_at
	SortOrder string // asc/desc
}

type Cat struct {
	ID uint64 `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`

	// ==================== 基础信息 ====================
	CatCode string `gorm:"type:varchar(64);uniqueIndex;not null;comment:猫咪编号" json:"cat_code"`
	Name    string `gorm:"type:varchar(64);not null;comment:猫咪名称" json:"name"`
	Breed   string `gorm:"type:varchar(64);comment:品种" json:"breed"`
	Color   string `gorm:"type:varchar(64);comment:毛色" json:"color"`
	Gender  string `gorm:"type:varchar(20);not null;default:'unknown';comment:性别 male/female/unknown" json:"gender"`

	// ==================== 体态信息 ====================
	BodySize  string     `gorm:"type:varchar(20);default:'medium';comment:体型 small/medium/large" json:"body_size"`
	AgeStage  string     `gorm:"type:varchar(30);default:'young';comment:年龄阶段 kitten/young/adult/old" json:"age_stage"`
	Weight    float32    `gorm:"type:decimal(5,2);default:0;comment:体重(kg)" json:"weight"`
	Character string     `gorm:"type:varchar(255);comment:性格描述" json:"character"`
	BirthDate *time.Time `gorm:"comment:预估出生日期" json:"birth_date"`

	// ==================== 展示信息 ====================
	Avatar      string `gorm:"type:varchar(500);comment:猫咪头像URL" json:"avatar"`
	Description string `gorm:"type:text;comment:简介" json:"description"`

	// ==================== 发现信息 ====================
	FoundAt          *time.Time `gorm:"comment:发现日期" json:"found_at"`
	DiscoveryAddress string     `gorm:"type:varchar(255);comment:发现地址" json:"discovery_address"`
	Longitude        float64    `gorm:"type:decimal(10,6);default:0;comment:经度" json:"longitude"`
	Latitude         float64    `gorm:"type:decimal(10,6);default:0;comment:纬度" json:"latitude"`

	// ==================== 健康信息 ====================
	IsVaccinated            bool       `gorm:"default:false;comment:是否已接种疫苗" json:"is_vaccinated"`
	IsHealthy               bool       `gorm:"default:false;comment:是否健康" json:"is_healthy"`
	NeedMedicalIntervention bool       `gorm:"default:true;comment:是否需要医疗干预" json:"need_medical_intervention"`
	SterilizationStatus     string     `gorm:"type:varchar(30);default:'unsterilized';index;comment:绝育状态 sterilized/unsterilized" json:"sterilization_status"`
	LastMedicalCheckAt      *time.Time `gorm:"comment:最近体检时间" json:"last_medical_check_at"`

	// ==================== 领养信息 ====================
	AdoptionStatus string     `gorm:"type:varchar(30);default:'pending';index;comment:领养状态 adopted/pending/unavailable" json:"adoption_status"`
	AdopterID      uint64     `gorm:"default:0;index;comment:领养人ID" json:"adopter_id"`
	AdoptedAt      *time.Time `gorm:"comment:领养时间" json:"adopted_at"`

	// ==================== 来源信息 ====================
	ApplyID   uint64 `gorm:"default:0;index;comment:来源申请单ID，0表示管理员直接创建" json:"apply_id"`
	CreatorID uint64 `gorm:"default:0;index;comment:档案创建人ID" json:"creator_id"`

	// ==================== 时间信息 ====================
	CreatedAt time.Time             `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time             `gorm:"comment:更新时间" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano;comment:删除时间" json:"deleted_at,omitempty"`
}

func (Cat) TableName() string {
	return "cats"
}

// ==================== 正式猫咪档案实现 ====================

func (r *repository) CreateCat(ctx context.Context, cat *Cat, tx ...*gorm.DB) error {
	setCatDefaults(cat)

	if err := validateCatForCreate(cat); err != nil {
		return err
	}

	return r.getDB(ctx, tx...).Create(cat).Error
}

func (r *repository) GetCatByID(ctx context.Context, catID uint64, tx ...*gorm.DB) (*Cat, error) {
	if catID == 0 {
		return nil, ErrInvalidCatID
	}

	var cat Cat
	err := r.getDB(ctx, tx...).
		Where("id = ?", catID).
		First(&cat).Error

	if err != nil {
		return nil, err
	}

	return &cat, nil
}

func (r *repository) GetCatByCode(ctx context.Context, catCode string, tx ...*gorm.DB) (*Cat, error) {
	if catCode == "" {
		return nil, ErrCatCodeRequired
	}

	var cat Cat
	err := r.getDB(ctx, tx...).
		Where("cat_code = ?", catCode).
		First(&cat).Error

	if err != nil {
		return nil, err
	}

	return &cat, nil
}

func (r *repository) ListCats(ctx context.Context, filter CatListFilter, tx ...*gorm.DB) ([]*Cat, int64, error) {
	db := r.getDB(ctx, tx...).Model(&Cat{})

	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		db = db.Where(
			"name LIKE ? OR cat_code LIKE ? OR breed LIKE ? OR color LIKE ?",
			keyword, keyword, keyword, keyword,
		)
	}

	if filter.CatCode != "" {
		db = db.Where("cat_code = ?", filter.CatCode)
	}

	if filter.Name != "" {
		db = db.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.Breed != "" {
		db = db.Where("breed = ?", filter.Breed)
	}

	if filter.Color != "" {
		db = db.Where("color = ?", filter.Color)
	}

	if filter.Gender != "" {
		db = db.Where("gender = ?", filter.Gender)
	}

	if filter.BodySize != "" {
		db = db.Where("body_size = ?", filter.BodySize)
	}

	if filter.AgeStage != "" {
		db = db.Where("age_stage = ?", filter.AgeStage)
	}

	if filter.SterilizationStatus != "" {
		db = db.Where("sterilization_status = ?", filter.SterilizationStatus)
	}

	if filter.AdoptionStatus != "" {
		db = db.Where("adoption_status = ?", filter.AdoptionStatus)
	}

	if filter.AdopterID > 0 {
		db = db.Where("adopter_id = ?", filter.AdopterID)
	}

	if filter.CreatorID > 0 {
		db = db.Where("creator_id = ?", filter.CreatorID)
	}

	if filter.ApplyID > 0 {
		db = db.Where("apply_id = ?", filter.ApplyID)
	}

	if filter.IsVaccinated != nil {
		db = db.Where("is_vaccinated = ?", *filter.IsVaccinated)
	}

	if filter.IsHealthy != nil {
		db = db.Where("is_healthy = ?", *filter.IsHealthy)
	}

	if filter.NeedMedicalIntervention != nil {
		db = db.Where("need_medical_intervention = ?", *filter.NeedMedicalIntervention)
	}

	if filter.FoundAtStart != nil {
		db = db.Where("found_at >= ?", *filter.FoundAtStart)
	}

	if filter.FoundAtEnd != nil {
		db = db.Where("found_at <= ?", *filter.FoundAtEnd)
	}

	if filter.CreatedAtStart != nil {
		db = db.Where("created_at >= ?", *filter.CreatedAtStart)
	}

	if filter.CreatedAtEnd != nil {
		db = db.Where("created_at <= ?", *filter.CreatedAtEnd)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit, offset := normalizePage(filter.Page, filter.PageSize)
	order := buildCatOrder(filter.SortBy, filter.SortOrder)

	var list []*Cat
	err := db.
		Order(order).
		Limit(limit).
		Offset(offset).
		Find(&list).Error

	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *repository) UpdateCat(ctx context.Context, catID uint64, values map[string]any, tx ...*gorm.DB) error {
	if catID == 0 {
		return ErrInvalidCatID
	}

	values = cleanCatUpdateValues(values)
	if len(values) == 0 {
		return nil
	}

	res := r.getDB(ctx, tx...).
		Model(&Cat{}).
		Where("id = ?", catID).
		Updates(values)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *repository) UpdateCatAdoptionStatus(
	ctx context.Context,
	catID uint64,
	status string,
	adopterID uint64,
	tx ...*gorm.DB,
) error {
	if catID == 0 {
		return ErrInvalidCatID
	}

	if !isValidAdoptionStatus(status) {
		return ErrInvalidAdoptionStatus
	}

	updates := map[string]any{
		"adoption_status": status,
	}

	switch status {
	case CatAdoptionStatusAdopted:
		if adopterID == 0 {
			return ErrAdopterIDRequired
		}

		now := time.Now()
		updates["adopter_id"] = adopterID
		updates["adopted_at"] = &now

	case CatAdoptionStatusPending:
		updates["adopter_id"] = 0
		updates["adopted_at"] = nil

	case CatAdoptionStatusUnavailable:
		updates["adopter_id"] = 0
		updates["adopted_at"] = nil
	}

	res := r.getDB(ctx, tx...).
		Model(&Cat{}).
		Where("id = ?", catID).
		Updates(updates)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *repository) DeleteCat(ctx context.Context, catID uint64, tx ...*gorm.DB) error {
	if catID == 0 {
		return ErrInvalidCatID
	}

	db := r.getDB(ctx, tx...)

	var cat Cat
	if err := db.Where("id = ?", catID).First(&cat).Error; err != nil {
		return err
	}

	// 建议：已领养猫咪不允许直接删除
	if cat.AdoptionStatus == CatAdoptionStatusAdopted {
		return ErrCannotDeleteAdoptedCat
	}

	res := db.Where("id = ?", catID).Delete(&Cat{})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// NearFilter 附近筛选参数
type NearFilter struct {
	Latitude  float64 // 纬度
	Longitude float64 // 经度
	Radius    float64 // 半径（公里）
}

// ListCatsWithNearby 支持附近查找的小猫列表（带距离）
func (r *repository) ListCatsWithNearby(ctx context.Context, filter CatListFilter, near *NearFilter, tx ...*gorm.DB) ([]*Cat, []float64, int64, error) {
	db := r.getDB(ctx, tx...).Model(&Cat{})

	// 应用基础筛选
	db = r.applyCatFilters(db, filter)

	// 只筛选有位置信息的猫咪
	if near != nil && near.Latitude != 0 && near.Longitude != 0 {
		db = db.Where("latitude != 0 AND longitude != 0")
	}

	// 先查询总数（不带距离计算）
	var total int64
	countDB := r.getDB(ctx, tx...).Model(&Cat{})
	countDB = r.applyCatFilters(countDB, filter)
	if near != nil && near.Latitude != 0 && near.Longitude != 0 {
		countDB = countDB.Where("latitude != 0 AND longitude != 0")
	}
	if err := countDB.Count(&total).Error; err != nil {
		return nil, nil, 0, err
	}

	// 分页
	limit, offset := normalizePage(filter.Page, filter.PageSize)
	order := buildCatOrder(filter.SortBy, filter.SortOrder)

	// 查询基础数据
	var cats []*Cat
	err := db.
		Order(order).
		Limit(limit).
		Offset(offset).
		Find(&cats).Error
	if err != nil {
		return nil, nil, 0, err
	}

	// 计算距离（如果有 near 参数）
	var distances []float64
	if near != nil && near.Latitude != 0 && near.Longitude != 0 {
		distances = make([]float64, len(cats))
		for i, cat := range cats {
			distances[i] = calculateDistance(
				near.Latitude, near.Longitude,
				cat.Latitude, cat.Longitude,
			)
		}

		// 如果需要按距离筛选（半径内）
		if near.Radius > 0 {
			filteredCats := make([]*Cat, 0)
			filteredDistances := make([]float64, 0)
			for i, cat := range cats {
				if distances[i] <= near.Radius {
					filteredCats = append(filteredCats, cat)
					filteredDistances = append(filteredDistances, distances[i])
				}
			}
			cats = filteredCats
			distances = filteredDistances
			total = int64(len(filteredCats)) // 重新计算总数
		}

		// 如果按距离排序
		if filter.SortBy == "distance" {
			// 按距离排序
			for i := 0; i < len(cats)-1; i++ {
				for j := i + 1; j < len(cats); j++ {
					if (filter.SortOrder == "asc" && distances[i] > distances[j]) ||
						(filter.SortOrder == "desc" && distances[i] < distances[j]) {
						cats[i], cats[j] = cats[j], cats[i]
						distances[i], distances[j] = distances[j], distances[i]
					}
				}
			}
		}
	}

	return cats, distances, total, nil
}

// applyCatFilters 应用筛选条件（提取公共逻辑）
func (r *repository) applyCatFilters(db *gorm.DB, filter CatListFilter) *gorm.DB {
	// 关键词搜索
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		db = db.Where(
			"name LIKE ? OR cat_code LIKE ? OR breed LIKE ? OR color LIKE ?",
			keyword, keyword, keyword, keyword,
		)
	}

	// 精确匹配
	if filter.CatCode != "" {
		db = db.Where("cat_code = ?", filter.CatCode)
	}
	if filter.Name != "" {
		db = db.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Breed != "" {
		db = db.Where("breed = ?", filter.Breed)
	}
	if filter.Color != "" {
		db = db.Where("color = ?", filter.Color)
	}
	if filter.Gender != "" {
		db = db.Where("gender = ?", filter.Gender)
	}

	// 体态信息
	if filter.BodySize != "" {
		db = db.Where("body_size = ?", filter.BodySize)
	}
	if filter.AgeStage != "" {
		db = db.Where("age_stage = ?", filter.AgeStage)
	}

	// 健康信息
	if filter.SterilizationStatus != "" {
		db = db.Where("sterilization_status = ?", filter.SterilizationStatus)
	}
	if filter.AdoptionStatus != "" {
		db = db.Where("adoption_status = ?", filter.AdoptionStatus)
	}
	if filter.IsVaccinated != nil {
		db = db.Where("is_vaccinated = ?", *filter.IsVaccinated)
	}
	if filter.IsHealthy != nil {
		db = db.Where("is_healthy = ?", *filter.IsHealthy)
	}
	if filter.NeedMedicalIntervention != nil {
		db = db.Where("need_medical_intervention = ?", *filter.NeedMedicalIntervention)
	}

	// 关联ID
	if filter.AdopterID > 0 {
		db = db.Where("adopter_id = ?", filter.AdopterID)
	}
	if filter.CreatorID > 0 {
		db = db.Where("creator_id = ?", filter.CreatorID)
	}
	if filter.ApplyID > 0 {
		db = db.Where("apply_id = ?", filter.ApplyID)
	}

	// 时间范围
	if filter.FoundAtStart != nil {
		db = db.Where("found_at >= ?", *filter.FoundAtStart)
	}
	if filter.FoundAtEnd != nil {
		db = db.Where("found_at <= ?", *filter.FoundAtEnd)
	}
	if filter.CreatedAtStart != nil {
		db = db.Where("created_at >= ?", *filter.CreatedAtStart)
	}
	if filter.CreatedAtEnd != nil {
		db = db.Where("created_at <= ?", *filter.CreatedAtEnd)
	}

	return db
}

// calculateDistance 计算两点距离（公里）- Haversine公式
func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371 // 地球半径（公里）

	// 转换为弧度
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLng := (lng2 - lng1) * math.Pi / 180

	// Haversine公式
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// GetAllCatsWithLocation 获取所有有坐标的猫咪（用于热点检测）
func (r *repository) GetAllCatsWithLocation(ctx context.Context, tx ...*gorm.DB) ([]*Cat, error) {
	db := r.getDB(ctx, tx...)
	var cats []*Cat

	err := db.WithContext(ctx).
		Where(`
	longitude != 0 
	AND latitude != 0 
	AND (
		need_medical_intervention = true 
		OR is_healthy = false
	)
`).
		Find(&cats).Error

	return cats, err
}

// GetCatsByIDs 批量获取猫咪信息
func (r *repository) GetCatsByIDs(ctx context.Context, ids []uint64, tx ...*gorm.DB) ([]*Cat, error) {
	if len(ids) == 0 {
		return []*Cat{}, nil
	}

	db := r.getDB(ctx, tx...)
	var cats []*Cat

	err := db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&cats).Error

	return cats, err
}

// GetCatsByHotspot 获取热点区域内的猫咪（用于任务自动创建）
func (r *repository) GetCatsByHotspot(ctx context.Context, centerLat, centerLng, radiusKm float64, tx ...*gorm.DB) ([]*Cat, error) {
	db := r.getDB(ctx, tx...)
	var cats []*Cat

	// 使用 MySQL 的空间函数计算距离
	// 粗略筛选：经纬度范围 ±radius度（1度≈111km）
	latRange := radiusKm / 111.0
	lngRange := radiusKm / (111.0 * math.Cos(centerLat*math.Pi/180))

	err := db.WithContext(ctx).
		Where("longitude BETWEEN ? AND ?", centerLng-lngRange, centerLng+lngRange).
		Where("latitude BETWEEN ? AND ?", centerLat-latRange, centerLat+latRange).
		Where("deleted_at IS NULL").
		Find(&cats).Error

	if err != nil {
		return nil, err
	}

	// 精确筛选距离
	result := make([]*Cat, 0)
	for _, c := range cats {
		dist := haversine(centerLat, centerLng, c.Latitude, c.Longitude)
		if dist <= radiusKm {
			result = append(result, c)
		}
	}

	return result, nil
}

// GetCatReportCounts 获取每只猫咪的近期报告/申请次数（用于优先级评分）
func (r *repository) GetCatReportCounts(ctx context.Context, catIDs []uint64, days int, tx ...*gorm.DB) (map[uint64]int, error) {
	db := r.getDB(ctx, tx...)

	type Result struct {
		CatID uint64
		Count int
	}

	var results []Result

	// 统计7天内的申请和任务创建次数
	threshold := time.Now().AddDate(0, 0, -days)

	err := db.WithContext(ctx).
		Table("cat_create_applies").
		Select("cat_id, COUNT(*) as count").
		Where("cat_id IN ?", catIDs).
		Where("created_at > ?", threshold).
		Group("cat_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	countMap := make(map[uint64]int)
	for _, r := range results {
		countMap[r.CatID] = r.Count
	}

	return countMap, nil
}

// todo 以下表结构暂时没有用到

// CatRescueRecord 救助历程表
type CatRescueRecord struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	CatID        uint64    `gorm:"index;not null;comment:猫咪ID" json:"cat_id"`
	RescueTime   time.Time `gorm:"comment:救助时间" json:"rescue_time"`
	RescuerName  string    `gorm:"type:varchar(64);comment:救助者姓名" json:"rescuer_name"`
	RescueStatus string    `gorm:"type:varchar(50);comment:救助状态" json:"rescue_status"`
	Description  string    `gorm:"type:text;comment:救助简介" json:"description"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (CatRescueRecord) TableName() string {
	return "cat_rescue_records"
}

type CatMedicalRecord struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	CatID        uint64    `gorm:"index;not null;comment:猫咪ID" json:"cat_id"`
	MedicalDate  time.Time `gorm:"comment:医疗日期" json:"medical_date"`
	MedicalType  string    `gorm:"type:varchar(50);comment:类型 vaccine/sterilization/check/treatment" json:"medical_type"`
	Content      string    `gorm:"type:text;comment:医疗内容" json:"content"`
	OperatorName string    `gorm:"type:varchar(64);comment:操作员" json:"operator_name"`

	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano;comment:删除时间" json:"deleted_at,omitempty"`
}

func (CatMedicalRecord) TableName() string {
	return "cat_medical_records"
}

/*
 * =========================================================
 * 修改猫咪档案申请表
 * =========================================================
 */

type CatUpdateApply struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	CatID uint64 `gorm:"index;not null;comment:猫咪ID" json:"cat_id"`

	Name        string `gorm:"type:varchar(64)" json:"name"`
	Gender      string `gorm:"type:varchar(20)" json:"gender"`
	BodySize    string `gorm:"type:varchar(20)" json:"body_size"`
	AgeStage    string `gorm:"type:varchar(30)" json:"age_stage"`
	Description string `gorm:"type:text" json:"description"`

	DiscoveryAddress string  `gorm:"type:varchar(255)" json:"discovery_address"`
	Longitude        float64 `gorm:"type:decimal(10,6);default:0" json:"longitude"`
	Latitude         float64 `gorm:"type:decimal(10,6);default:0" json:"latitude"`

	ImageURLs string `gorm:"type:json;comment:图片JSON" json:"image_urls"`

	ChangeReason string `gorm:"type:text;comment:修改原因" json:"change_reason"`

	ApplicantUserID uint64 `gorm:"index;not null" json:"applicant_user_id"`

	Status       string `gorm:"type:varchar(30);default:'pending';index" json:"status"`
	RejectReason string `gorm:"type:text" json:"reject_reason"`
	ReviewerID   uint64 `gorm:"default:0" json:"reviewer_id"`

	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (CatUpdateApply) TableName() string {
	return "cat_update_applies"
}

/*
 * =========================================================
 * 医疗记录申请表
 * =========================================================
 */

type CatMedicalApply struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	CatID uint64 `gorm:"index;not null" json:"cat_id"`

	MedicalType string `gorm:"type:varchar(50);comment:vaccine/sterilization/check/treatment" json:"medical_type"`
	Content     string `gorm:"type:text" json:"content"`

	ApplicantUserID uint64 `gorm:"index;not null" json:"applicant_user_id"`

	Status       string `gorm:"type:varchar(30);default:'pending';index" json:"status"`
	RejectReason string `gorm:"type:text" json:"reject_reason"`
	ReviewerID   uint64 `gorm:"default:0" json:"reviewer_id"`

	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (CatMedicalApply) TableName() string {
	return "cat_medical_applies"
}

/*
 * =========================================================
 * 救助记录申请表
 * =========================================================
 */

type CatRescueApply struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	CatID uint64 `gorm:"index;not null" json:"cat_id"`

	RescueTime   time.Time `json:"rescue_time"`
	RescuerName  string    `gorm:"type:varchar(64)" json:"rescuer_name"`
	RescueStatus string    `gorm:"type:varchar(50)" json:"rescue_status"`
	Description  string    `gorm:"type:text" json:"description"`

	ApplicantUserID uint64 `gorm:"index;not null" json:"applicant_user_id"`

	Status       string `gorm:"type:varchar(30);default:'pending';index" json:"status"`
	RejectReason string `gorm:"type:text" json:"reject_reason"`
	ReviewerID   uint64 `gorm:"default:0" json:"reviewer_id"`

	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (CatRescueApply) TableName() string {
	return "cat_rescue_applies"
}
