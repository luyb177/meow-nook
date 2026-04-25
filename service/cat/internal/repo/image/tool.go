package image

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// getDB 根据是否传入 tx 返回对应的 DB 实例
func (r *repository) getDB(ctx context.Context, tx ...*gorm.DB) *gorm.DB {
	db := r.db.WithContext(ctx)
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].WithContext(ctx)
	}
	return db
}

// withTx 工具函数：封装事务逻辑
// 如果外部传入 tx，则在外部事务中执行 fn。
// 如果外部未传入 tx，则开启一个新的内部事务执行 fn。
// 这样可以实现内部事务的“智能”管理：可独立运行，也可加入外部事务。
func (r *repository) withTx(ctx context.Context, fn func(tx *gorm.DB) error, txs ...*gorm.DB) error {
	if len(txs) > 0 && txs[0] != nil {
		// 外部提供了事务，在外部事务中执行
		return fn(txs[0].WithContext(ctx))
	}

	// 外部未提供事务，开启一个新事务执行
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

func (r *repository) validateTarget(targetType string, targetID uint64) error {
	if targetType == "" || targetID == 0 {
		return ErrInvalidTarget
	}
	return nil
}

func (r *repository) validateImage(img *Image) error {
	if img == nil {
		return errors.New("image cannot be nil for validation")
	}
	if err := r.validateTarget(img.TargetType, img.TargetID); err != nil {
		return err
	}
	if img.URL == "" {
		return ErrEmptyURL
	}
	return nil
}

func (r *repository) validateImages(images []*Image) error {
	coverCount := make(map[string]int) // Key: targetType_targetID

	for _, img := range images {
		if err := r.validateImage(img); err != nil {
			return err
		}

		if img.IsCover {
			key := fmt.Sprintf("%s_%d", img.TargetType, img.TargetID)
			coverCount[key]++
			if coverCount[key] > 1 {
				return fmt.Errorf("%w: target_type: %s, target_id: %d", ErrMultiCover, img.TargetType, img.TargetID)
			}
		}
	}
	return nil
}
