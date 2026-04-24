package verify

import (
	_ "embed"

	"github.com/redis/go-redis/v9"
)

//go:embed lua/get-and-delete.lua
var getAndDeleteLua string

var getAndDeleteScript = redis.NewScript(getAndDeleteLua)
