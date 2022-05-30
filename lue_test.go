package main

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

func TestIncrLimit(t *testing.T) {
	ctx := context.Background()
	_ = rdb.FlushDB(ctx).Err()

	var incrLimit = incrLue()
	keys := []string{"hello"}
	limit := []interface{}{10}

	thread := 20
	//var wg sync.WaitGroup
	ch := make(chan bool, thread)
	for i := 1; i <= thread; i++ {
		//wg.Add(1)
		go func(i int) {
			//defer wg.Done()
			_, err := incrLimit.Run(ctx, rdb, keys, limit).Int()
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
	//wg.Wait()
	expected := 10
	if count != expected {
		t.Errorf("expect %v got %v\n", expected, count)
	}
}
