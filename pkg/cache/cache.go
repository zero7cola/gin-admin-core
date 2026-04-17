package cache

import (
	"encoding/json"
	"github.com/spf13/cast"
	"github.com/zero7cola/gin-admin-core/pkg/logger"
	"sync"
	"time"
)

type CacheService struct {
	Store Store
}

var once sync.Once
var Cache *CacheService

func InitWithCacheStore(store Store) {
	once.Do(func() {
		Cache = &CacheService{
			Store: store,
		}
	})
}

func Set(key string, obj interface{}, expireTime time.Duration) {
	b, err := json.Marshal(&obj)
	logger.LogIf(err)
	Cache.Store.Set(key, string(b), expireTime)
}

func Get(key string) interface{} {
	stringValue := Cache.Store.Get(key)
	var wanted interface{}
	err := json.Unmarshal([]byte(stringValue), &wanted)
	logger.LogIf(err)
	return wanted
}

func Has(key string) bool {
	return Cache.Store.Has(key)
}

// GetObject 应该传地址，用法如下:
//
//	model := user.User{}
//	cache.GetObject("key", &model)
func GetObject(key string, wanted interface{}) {
	val := Cache.Store.Get(key)
	if len(val) > 0 {
		err := json.Unmarshal([]byte(val), &wanted)
		logger.LogIf(err)
	}
}

func GetString(key string) string {
	return cast.ToString(Get(key))
}

func Forget(key string) {
	Cache.Store.Forget(key)
}

func Forever(key string, value string) {
	Cache.Store.Set(key, value, 0)
}

func Flush() {
	Cache.Store.Flush()
}

func Increment(parameters ...interface{}) {
	Cache.Store.Increment(parameters...)
}

func Decrement(parameters ...interface{}) {
	Cache.Store.Decrement(parameters...)
}

func IsAlive() error {
	return Cache.Store.IsAlive()
}
