package adoption

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
