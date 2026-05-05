package tag

import (
	"context"

	"gorm.io/gorm"
)

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

func normalizePage(page, pageSize int) (limit, offset int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return pageSize, (page - 1) * pageSize
}
