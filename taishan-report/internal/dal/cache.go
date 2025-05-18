package dal

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	localCache *cache.Cache
)

const (
	DefaultExpirationDuration = 5 * time.Minute
)

func MustInitLocalCache() {
	localCache = cache.New(10*time.Minute, 2*time.Minute)
}

func LocalCacheSet(key string, data any) {
	localCache.Set(key, data, DefaultExpirationDuration)
}

func LocalCacheGet(key string) (res any, success bool) {
	res, success = localCache.Get(key)
	return
}
