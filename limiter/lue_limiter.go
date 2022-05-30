package limiter

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type lueLimiter struct {
	*redis.Client
	*redis.Script
}

func NewLueLimiter(rdb *redis.Client) *lueLimiter {
	incr := redis.NewScript(`
		local key = KEYS[1]
		local limit = tonumber(ARGV[1])
		redis.log(redis.LOG_NOTICE, "key", key, "limit", limit)
		local value = redis.call("incr", key)
		if( value > limit )
		then
			return redis.error_reply('over limit')
		end
		return value
	`)
	return &lueLimiter{Client: rdb, Script: incr}
}

func (l *lueLimiter) IncrWithLimit(key string, limit int) (int, error) {
	ctx := context.Background()
	keys := []string{key}
	args := []interface{}{limit}
	num, err := l.Script.Run(ctx, l.Client, keys, args).Int()
	return num, err
}
