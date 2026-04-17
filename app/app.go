package app

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zero7cola/gin-admin-core/config"
	"github.com/zero7cola/gin-admin-core/core"
	"github.com/zero7cola/gin-admin-core/pkg/database"
	"gorm.io/gorm"
)

type App struct {
	engine *gin.Engine
	group  *gin.RouterGroup
	prefix string

	middlewares []gin.HandlerFunc
	modules     []core.Module
	db          *gorm.DB
	redis       *redis.Client
}

func (a *App) DB() *gorm.DB {
	return a.db
}

func (a *App) Redis() *redis.Client {
	return a.redis
}

func (a *App) Group(path string) *gin.RouterGroup {
	return a.group.Group(path)
}

func (a *App) Use(m ...gin.HandlerFunc) {
	a.middlewares = append(a.middlewares, m...)
}

func (a *App) Register(m core.Module) {
	a.modules = append(a.modules, m)
}

type Option func(*App)

func WithSetup(fn func(app *App)) Option {
	return func(a *App) {
		fn(a)
	}
}

func WithPrefix(prefix string) Option {
	return func(a *App) {
		a.prefix = prefix
	}
}

func registerBuiltinModules(app *App) {
}

func initDB(cfg *config.Config) *gorm.DB {
	return database.Init(nil, nil)
}

func newApp(r *gin.Engine, prefix string, cfg *config.Config) *App {
	return &App{
		engine: r,
		prefix: prefix,
	}
}

func Init(r *gin.Engine, cfg *config.Config, opts ...Option) {
	app := newApp(r, "/admin", cfg)

	// 初始化基础设施
	app.db = initDB(cfg)

	// 执行外部配置
	for _, opt := range opts {
		opt(app)
	}

	// 加载内部modules
	registerBuiltinModules(app)

	// 创建路由组
	app.group = r.Group(app.prefix)

	// 注册中间件
	if len(app.middlewares) > 0 {
		app.group.Use(app.middlewares...)
	}

	// 注册外部模块
	for _, m := range app.modules {
		m.Register(r)
	}
}

func IsLocal() bool {
	return config.Get("app.env") == "local"
}

func IsProduction() bool {
	return config.Get("app.env") == "production"
}

func IsTesting() bool {
	return config.Get("app.env") == "testing"
}

func IsDebug() bool {
	return config.GetBool("debug") == true
}

func GetFileUrl(fileName string) string {
	return config.GetString("app.file_url") + "/" + fileName
}

// TimenowInTimezone 获取当前时间，支持时区
func TimenowInTimezone() time.Time {
	chinaTimezone, _ := time.LoadLocation(config.GetString("app.timezone"))
	return time.Now().In(chinaTimezone)
}
