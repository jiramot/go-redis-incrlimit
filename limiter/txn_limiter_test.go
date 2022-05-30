package limiter

import (
	"context"
	"github.com/go-redis/redis/v8"
	"testing"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
}

func TestIncrLimitPipeline(t *testing.T) {
	ctx := context.Background()
	_ = rdb.FlushDB(ctx).Err()

	limit := TxnLimiter(rdb)
	thread := 20
	ch := make(chan bool, thread)
	for i := 1; i <= thread; i++ {
		go func(i int) {
			_, err := limit.IncrWithLimit("key", 10)
			if err != nil {
				ch <- false
			} else {
				ch <- true
			}
		}(i)
	}
	count := 0
	for i := 1; i <= thread; i++ {
		rs := <-ch
		if rs {
			count++
		}

	}
	expected := 10
	if count != expected {
		t.Errorf("expect %v got %v\n", expected, count)
	}
}
