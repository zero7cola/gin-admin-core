package cache

import (
	"fmt"
	"github.com/zero7cola/gin-admin-core/config"
	redis2 "github.com/zero7cola/gin-admin-core/pkg/redis"
	"time"
)

type RedisStore struct {
	RedisClient *redis2.RedisClient
	KeyPrefix   string
}

func NewRedisStore(addr, username, password string, db int) *RedisStore {
	return &RedisStore{
		RedisClient: redis2.NewClient(addr, username, password, db),
		KeyPrefix:   fmt.Sprintf("%s:cache:", config.GetString("app_name")),
	}
}

func (s *RedisStore) Set(key string, value string, timeout time.Duration) {
	s.RedisClient.Set(s.KeyPrefix+key, value, timeout)
}

func (s *RedisStore) Get(key string) string {
	return s.RedisClient.Get(s.KeyPrefix + key)
}

func (s *RedisStore) Has(key string) bool {
	return s.RedisClient.Has(s.KeyPrefix + key)
}

func (s *RedisStore) Forget(key string) {
	s.RedisClient.Del(s.KeyPrefix + key)
}

func (s *RedisStore) Forever(key string, value string) {
	s.RedisClient.Set(s.KeyPrefix+key, value, 0)
}

func (s *RedisStore) Flush() {
	s.RedisClient.FlushDB()
}

func (s *RedisStore) Increment(parameters ...interface{}) {
	s.RedisClient.Increment(parameters...)
}

func (s *RedisStore) Decrement(parameters ...interface{}) {
	s.RedisClient.Decrement(parameters...)
}

func (s *RedisStore) IsAlive() error {
	return s.RedisClient.Ping()
}
