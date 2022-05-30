package main

import (
	"context"
	"github.com/jiramot/go-redis-incrlimit/limiter"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	_ = rdb.FlushDB(ctx).Err()

	lue := limiter.NewLueLimiter(rdb)
	txn := limiter.TxnLimiter(rdb)

	val, _ := lue.IncrWithLimit("key", 2)
	println(val)
	val, _ = txn.IncrWithLimit("key", 2)
	println(val)
	val, _ = lue.IncrWithLimit("key", 2)
	println(val)
	val, _ = txn.IncrWithLimit("key", 2)
	println(val)
}
