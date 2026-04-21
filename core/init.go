package core

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/zero7cola/gin-admin-core/pkg/cache"
	"github.com/zero7cola/gin-admin-core/pkg/database"
	"github.com/zero7cola/gin-admin-core/pkg/logger"
	redisClient "github.com/zero7cola/gin-admin-core/pkg/redis"
	"github.com/zero7cola/gin-admin-core/setting"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type InitConfig struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Logger *zap.Logger
	Config *setting.Setting
	Cache  *cache.CacheService
}

type Option func(config *InitConfig)

func loadConfig(path string) (*InitConfig, error) {
	v := viper.New()

	if len(path) > 0 {
		v.SetConfigType("yaml") // 类型
		v.AddConfigPath(".")    // 当前目录
		v.SetConfigFile(path)

		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	var cfg InitConfig
	if err := v.Unmarshal(&cfg.Config); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func WithJWTConfig(jwt setting.JWTConfig) Option {
	return func(c *InitConfig) {
		c.Config.JWT = jwt
	}
}

func WithAppConfig(app setting.AppConfig) Option {
	return func(c *InitConfig) {
		c.Config.App = app
	}
}

func WithStorageConfig(storage setting.StorageConfig) Option {
	return func(c *InitConfig) {
		c.Config.Storage = storage
	}
}

func WithCaptchaConfig(config setting.CaptchaConfig) Option {
	return func(c *InitConfig) {
		c.Config.Captcha = config
	}
}

func WithPagingConfig(config setting.PagingConfig) Option {
	return func(c *InitConfig) {
		c.Config.Paging = config
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

func WithCache(cache *cache.CacheService) Option {
	return func(c *InitConfig) {
		c.Cache = cache
	}
}

func InitWithFile(path string, opts ...Option) error {
	// 1️⃣ 先加载文件
	cfg, err := loadConfig(path)
	if err != nil {
		return err
	}

	// 2️⃣ 应用手动配置（覆盖）
	for _, opt := range opts {
		opt(cfg)
	}

	// 3️⃣ 初始化内部 Context
	internalInit(cfg)

	return nil
}

func Init(opts ...Option) error {
	// 1️⃣ 先加载文件
	cfg := &InitConfig{}

	// 2️⃣ 应用手动配置（覆盖）
	for _, opt := range opts {
		opt(cfg)
	}

	// 3️⃣ 初始化内部 Context
	internalInit(cfg)

	return nil
}

func internalInit(c *InitConfig) {
	if c == nil {
		panic("setting is nil")
	}

	if c.Logger == nil {
		panic("logger is required")
	}

	if c.DB == nil {
		panic("db is required")
	}

	logger.Logger = c.Logger
	database.DB = c.DB
	setting.GlobalSetting = c.Config

	if c.Redis == nil {
		panic("redis is required")
	} else {
		rClient := redisClient.NewClient(c.Redis)
		redisClient.Redis = rClient
	}

	if c.Cache != nil {
		cache.Cache = c.Cache
	}

}
