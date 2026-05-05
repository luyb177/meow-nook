package task

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"
)

// getDB 获取数据库连接（支持事务）
func (r *repository) getDB(ctx context.Context, tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		return tx[0]
	}
	return r.db.WithContext(ctx)
}

// ==================== 辅助函数 ====================

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return pageSize, (page - 1) * pageSize
}

func ptr[T any](v T) *T {
	return &v
}

func structToJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
