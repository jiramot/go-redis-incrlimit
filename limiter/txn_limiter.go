package limiter

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

type txnLimiter struct {
	*redis.Client
}

func TxnLimiter(rdb *redis.Client) *txnLimiter {
	return &txnLimiter{Client: rdb}
}

func (l *txnLimiter) IncrWithLimit(key string, limit int) (int, error) {
	ctx := context.Background()
	var v int
	txf := func(tx *redis.Tx) error {
		val, err := tx.Incr(ctx, key).Uint64()
		if err != nil {
			return err
		}
		v = int(val)
		if v > limit {
			return errors.New("over limit")
		}
		return nil
	}

	err := l.Client.Watch(ctx, txf, key)
	if err != nil {
		return 0, err
	}
	return v, nil
}
