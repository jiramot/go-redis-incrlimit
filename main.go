package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func incrLue() *redis.Script {
	return redis.NewScript(`
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
}

func main() {
	fmt.Println("start")
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	_ = rdb.FlushDB(ctx).Err()

	var incrLimit = incrLue()
	keys := []string{"hello"}
	limit := []interface{}{10}

	for i := 1; i < 20; i++ {
		fmt.Printf("loop %d\n", i)
		num, err := incrLimit.Run(ctx, rdb, keys, limit).Int()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(num)
	}
}
