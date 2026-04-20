package core

import (
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/redis/go-redis/v9"
	redisClient "github.com/zero7cola/gin-admin-core/pkg/redis"
	"go.uber.org/zap"
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
			log.Fatalf("Error reading config file, %s", err)
		}
	}

	var cfg InitConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type Option func(config *InitConfig)

func WithJWTSecret(jwt JWTConfig) Option {
	return func(c *InitConfig) {
		c.Config.JWT = jwt
	}
}

func WithApp(app AppConfig) Option {
	return func(c *InitConfig) {
		c.Config.App = app
	}
}

func WithStorage(storage StorageConfig) Option {
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

type InitConfig struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Logger *zap.Logger
	Config *Config
}

type appContext struct {
	DB     *gorm.DB
	Redis  *redisClient.RedisClient
	Logger *zap.Logger
	Config *Config
}

type Config struct {
	App     AppConfig
	Storage StorageConfig
	JWT     JWTConfig
	Captcha CaptchaConfig
	Paging  PagingConfig
}

type AppConfig struct {
	Name     string `mapstructure:"name" yaml:"name"`
	Key      string `mapstructure:"key" yaml:"key"`
	Url      string `mapstructure:"url" yaml:"url"`
	HttpPort string `mapstructure:"http_port" yaml:"http_port"`
	FileUrl  string `mapstructure:"file_url" yaml:"file_url"`
	Env      string `mapstructure:"env" yaml:"env"`
	Version  string `mapstructure:"version" yaml:"version"`
	Debug    bool   `mapstructure:"debug" yaml:"debug"`
	Timezone string `mapstructure:"timezone" yaml:"timezone"`
}

type StorageConfig struct {
	Driver    string              `mapstructure:"driver" yaml:"driver"`
	SizeLimit int64               `mapstructure:"size_limit" yaml:"size_limit"`
	Ext       []string            `mapstructure:"ext" yaml:"ext"`
	Local     *LocalStorageConfig `mapstructure:"local" yaml:"local"`
	Oss       *OssStorageConfig   `mapstructure:"oss" yaml:"oss"`
}

type LocalStorageConfig struct {
	Path         string `mapstructure:"path" yaml:"path"`
	Domain       string `mapstructure:"domain" yaml:"domain"`
	StaticPrefix string `mapstructure:"static" yaml:"static"`
}

type OssStorageConfig struct {
	KeyId     string `mapstructure:"key_id" yaml:"key_id"`
	KeySecret string `mapstructure:"key_secret" yaml:"key_secret"`
	Region    string `mapstructure:"region" yaml:"region"`
	Bucket    string `mapstructure:"bucket" yaml:"bucket"`
	Domain    string `mapstructure:"domain" yaml:"domain"`
}

type PagingConfig struct {
	PerPage         int    `mapstructure:"perpage" yaml:"perpage"`
	UrlQueryOrder   string `mapstructure:"url_query_order" yaml:"url_query_order"`
	UrlQuerySort    string `mapstructure:"url_query_sort" yaml:"url_query_sort"`
	UrlQueryPage    string `mapstructure:"url_query_page" yaml:"url_query_page"`
	UrlQueryPerPage string `mapstructure:"url_query_per_page" yaml:"url_query_per_page"`
}

type CaptchaConfig struct {
	Height     int     `mapstructure:"height" yaml:"height"`
	Width      int     `mapstructure:"width" yaml:"width"`
	Length     int     `mapstructure:"length" yaml:"length"`
	Maxskew    float64 `mapstructure:"maxskew" yaml:"maxskew"`
	Dotcount   int     `mapstructure:"dotcount" yaml:"dotcount"`
	ExpireTime int     `mapstructure:"expire_time" yaml:"expire_time"`
}

type JWTConfig struct {
	ExpireTime     int `mapstructure:"expire_time" yaml:"expire_time"`           // 过期时间，单位是分钟，一般不超过两个小时
	MaxReFreshTime int `mapstructure:"max_refresh_time" yaml:"max_refresh_time"` // 允许刷新时间，单位分钟，从 Token 的签名时间算起
}

var Global *appContext

func Init(c *InitConfig) {
	if c == nil {
		panic("config is nil")
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
		Config: c.Config,
	}

	if c.Redis == nil {
		panic("redis is required")
	} else {
		Global.Redis = redisClient.NewClient(c.Redis)
	}

}

func IsLocal() bool {
	return Global.Config.App.Env == "local"
}

func IsProduction() bool {
	return Global.Config.App.Env == "production"
}

func IsTesting() bool {
	return Global.Config.App.Env == "testing"
}

func IsDebug() bool {
	return Global.Config.App.Debug == true
}

// TimenowInTimezone 获取当前时间，支持时区
func TimenowInTimezone() time.Time {
	chinaTimezone, _ := time.LoadLocation(Global.Config.App.Timezone)
	return time.Now().In(chinaTimezone)
}

// URL 传参 path 拼接站点的 URL
func URL(path string) string {
	return Global.Config.App.Url + "/" + path
}

// VADMINURL 拼接带 admin 标示 URL
func VADMINURL(path string) string {
	return URL("/admin/" + path)
}

// V1URL 拼接带 v1 标示 URL
func V1URL(path string) string {
	return URL("/v1/" + path)
}
