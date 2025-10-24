package cache

import "time"

type Cache interface {
	Set(key string, value string, ttl time.Duration)
	Get(key string) (string, bool)
	Close()
}
