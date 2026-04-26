package repo

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/image"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/sequence"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/tag"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repositories struct {
	Cat      cat.Repository
	Image    image.Repository
	Tag      tag.Repository
	Sequence sequence.Repository
	db       *gorm.DB
}

func NewRepository(db *gorm.DB, client *redis.Client) *Repositories {
	return &Repositories{
		Cat:      cat.NewRepository(db, client),
		Image:    image.NewRepository(db, client),
		Tag:      tag.NewRepository(db, client),
		Sequence: sequence.NewRepository(db, client),
		db:       db,
	}
}

// WithTx 开启事务，上层 logic 直接调用这个函数
// fn 里面的 repo 操作都使用同一个 tx
func (r *Repositories) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}
