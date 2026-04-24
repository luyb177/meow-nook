package repo

import (
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repositories struct {
	Cat cat.Repository
}

func NewRepository(db *gorm.DB, client *redis.Client) *Repositories {
	return &Repositories{
		Cat: cat.NewRepository(db, client),
	}
}
