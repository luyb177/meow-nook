package cat

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ---------- 工具函数 ----------

func (r *repository) getDB(ctx context.Context, tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *repository) withTx(ctx context.Context, fn func(tx *gorm.DB) error, tx ...*gorm.DB) error {
	if len(tx) > 0 && tx[0] != nil {
		return fn(tx[0].WithContext(ctx))
	}
	return r.db.WithContext(ctx).Transaction(fn)
}

func normalizePage(page, pageSize int) (limit int, offset int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	limit = pageSize
	offset = (page - 1) * pageSize
	return
}

func setCatDefaults(cat *Cat) {
	if cat.Gender == "" {
		cat.Gender = CatGenderUnknown
	}
	if cat.BodySize == "" {
		cat.BodySize = CatBodySizeMedium
	}
	if cat.AgeStage == "" {
		cat.AgeStage = CatAgeStageYoung
	}
	if cat.SterilizationStatus == "" {
		cat.SterilizationStatus = CatSterilizationStatusUnsterilized
	}
	if cat.AdoptionStatus == "" {
		cat.AdoptionStatus = CatAdoptionStatusPending
	}

	// 默认健康状态按你的业务来：
	// 你目前结构体里 IsHealthy=false, NeedMedicalIntervention=true，
	// 表示刚入库默认需要检查/医疗干预，这个可以保留。
}

func validateCatForCreate(cat *Cat) error {
	if cat == nil {
		return errors.New("cat is nil")
	}
	if cat.CatCode == "" {
		return ErrCatCodeRequired
	}
	if cat.Name == "" {
		return ErrCatNameRequired
	}

	if cat.AdoptionStatus == CatAdoptionStatusAdopted && cat.AdopterID == 0 {
		return ErrAdopterIDRequired
	}

	if cat.AdoptionStatus == CatAdoptionStatusAdopted && cat.AdoptedAt == nil {
		now := time.Now()
		cat.AdoptedAt = &now
	}

	return nil
}

func isValidAdoptionStatus(status string) bool {
	switch status {
	case CatAdoptionStatusPending,
		CatAdoptionStatusAdopted,
		CatAdoptionStatusUnavailable:
		return true
	default:
		return false
	}
}

func buildCatOrder(sortBy, sortOrder string) string {
	allowedSortBy := map[string]string{
		"id":         "id",
		"created_at": "created_at",
		"updated_at": "updated_at",
		"found_at":   "found_at",
	}

	column, ok := allowedSortBy[sortBy]
	if !ok {
		column = "id"
	}

	order := "DESC"
	if strings.ToLower(sortOrder) == "asc" {
		order = "ASC"
	}

	return column + " " + order
}

func cleanCatUpdateValues(values map[string]any) map[string]any {
	if len(values) == 0 {
		return values
	}

	// 这些字段不建议通过通用 UpdateCat 修改
	blockedFields := map[string]struct{}{
		"id":         {},
		"cat_code":   {}, // 如果需要修改编号，建议单独做接口
		"created_at": {},
		"updated_at": {},
		"deleted_at": {},

		// 领养状态建议走 UpdateCatAdoptionStatus
		"adoption_status": {},
		"adopter_id":      {},
		"adopted_at":      {},
	}

	for field := range blockedFields {
		delete(values, field)
	}

	return values
}
