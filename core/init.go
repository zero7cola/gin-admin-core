package core

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	redisClient "github.com/zero7cola/gin-admin-core/pkg/redis"
	"github.com/zero7cola/gin-admin-core/setting"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
)

func LoadConfig(path string) (*InitConfig, error) {
	v := viper.New()

	if len(path) > 0 {
		v.SetConfigFile(path)

		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
		viper.WatchConfig()
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading setting file, %s", err)
		}
	}

	var cfg InitConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func WithJWTSecret(jwt setting.JWTConfig) Option {
	return func(c *InitConfig) {
		c.Config.JWT = jwt
	}
}

func WithApp(app setting.AppConfig) Option {
	return func(c *InitConfig) {
		c.Config.App = app
	}
}

func WithStorage(storage setting.StorageConfig) Option {
	return func(c *InitConfig) {
		c.Config.Storage = storage
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(c *InitConfig) {
		c.Logger = logger
	}
}

func WithDB(db *gorm.DB) Option {
	return func(c *InitConfig) {
		c.DB = db
	}
}

func WithRedis(redis *redis.Client) Option {
	return func(c *InitConfig) {
		c.Redis = redis
	}
}

type InitConfig struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Logger *zap.Logger
	Config *setting.Setting
}

type appContext struct {
	DB     *gorm.DB
	Redis  *redisClient.RedisClient
	Logger *zap.Logger
}

var Global *appContext

type Option func(config *InitConfig)

func InitWithFile(path string, opts ...Option) error {
	// 1️⃣ 先加载文件
	cfg, err := LoadConfig(path)
	if err != nil {
		return err
	}

	// 2️⃣ 应用手动配置（覆盖）
	for _, opt := range opts {
		opt(cfg)
	}

	// 3️⃣ 初始化内部 Context
	Init(cfg)

	return nil
}

func Init(c *InitConfig) {
	if c == nil {
		panic("setting is nil")
	}

	if c.Logger == nil {
		panic("logger is required")
	}

	if c.DB == nil {
		panic("db is required")
	}

	Global = &appContext{
		DB:     c.DB,
		Logger: c.Logger,
	}

	setting.GlobalSetting = c.Config

	if c.Redis == nil {
		panic("redis is required")
	} else {
		Global.Redis = redisClient.NewClient(c.Redis)
	}

}
