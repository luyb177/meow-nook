package repo

import (
	"github.com/luyb177/meow-nook/service/user/internal/repo/verify"
	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	Verify verify.Repository
}

func NewRepositories(redisClient *redis.Client) *Repositories {
	return &Repositories{
		Verify: verify.NewVerifyRepo(redisClient),
	}
}
