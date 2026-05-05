package repo

import (
	userrepo "github.com/luyb177/meow-nook/service/user/internal/repo/user"
	"github.com/luyb177/meow-nook/service/user/internal/repo/verify"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repositories struct {
	Verify verify.Repository
	User   userrepo.Repository
}

func NewRepositories(redisClient *redis.Client) *Repositories {
	return &Repositories{
		Verify: verify.NewVerifyRepo(redisClient),
		// User repo 需要 MySQL，由 NewRepositoriesWithDB 初始化
	}
}

// NewRepositoriesWithDB creates Repositories with both Redis and MySQL.
func NewRepositoriesWithDB(redisClient *redis.Client, db *gorm.DB) *Repositories {
	return &Repositories{
		Verify: verify.NewVerifyRepo(redisClient),
		User:   userrepo.NewUserRepo(db),
	}
}
