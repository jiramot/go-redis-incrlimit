package limiter

type Limiter interface {
	IncrWithLimit(key string, limit int) (int, error)
}
