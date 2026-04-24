package cat

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface{}

type repo struct {
	db     *gorm.DB
	client *redis.Client
}

func NewRepository(db *gorm.DB, client *redis.Client) Repository {
	return &repo{
		db:     db,
		client: client,
	}
}
