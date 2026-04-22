package core

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/zero7cola/gin-admin-core/model"
	"github.com/zero7cola/gin-admin-core/model/adminMenu"
	"github.com/zero7cola/gin-admin-core/model/adminOperationLog"
	"github.com/zero7cola/gin-admin-core/model/adminPermission"
	"github.com/zero7cola/gin-admin-core/model/adminRole"
	"github.com/zero7cola/gin-admin-core/model/adminUser"
	configModel "github.com/zero7cola/gin-admin-core/model/config"
	fileModel "github.com/zero7cola/gin-admin-core/model/file"
	"github.com/zero7cola/gin-admin-core/pkg/cache"
	"github.com/zero7cola/gin-admin-core/pkg/database"
	"github.com/zero7cola/gin-admin-core/pkg/logger"
	redisClient "github.com/zero7cola/gin-admin-core/pkg/redis"
	"github.com/zero7cola/gin-admin-core/setting"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	//
	err := insertInitData()

	if err != nil {
		logger.LogIf(err)
		panic(err)
	}

}

func insertInitData() error {

	err := database.DB.AutoMigrate(
		&adminUser.AdminUser{},
		&adminRole.AdminRole{},
		&adminMenu.AdminMenu{},
		&adminPermission.AdminPermission{},
		&fileModel.File{},
		&adminOperationLog.AdminOperationLog{},
		&configModel.Config{},
	)

	if err != nil {
		return err
	}

	err = seedAdminUser()

	if err != nil {
		return err
	}

	err = seedAdminMenus()

	if err != nil {
		return err
	}

	return nil
}

func insertIgnoreOrUpdate(data interface{}, isUp bool) error {
	if isUp {
		return database.DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(data).Error
	}

	return database.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(data).Error
}

func seedAdminUser() error {
	var users = []adminUser.AdminUser{
		{
			BaseModel: model.BaseModel{
				ID: 1,
			},
			Username: "admin",
			Password: "$2a$14$UPDOeuhOq6k6o2jnp3rCnudpcogjfSImV9hsHjKSEuMsPdoWY9Pk6",
			Name:     "Administrator",
		},
	}

	return insertIgnoreOrUpdate(users, false)
}

func seedAdminMenus() error {
	var menus = []adminMenu.AdminMenu{
		{
			BaseModel: model.BaseModel{
				ID: 1,
			},
			ParentId: 0,
			Order:    1,
			Name:     "管理后台",
			Icon:     "el-icon-s-management",
			Uri:      "",
		},
		{
			BaseModel: model.BaseModel{
				ID: 2,
			},
			ParentId: 1,
			Order:    2,
			Name:     "管理员",
			Icon:     "el-icon-user",
			Uri:      "admin/users/index",
		},
		{
			BaseModel: model.BaseModel{
				ID: 3,
			},
			ParentId: 1,
			Order:    3,
			Name:     "角色",
			Icon:     "el-icon-user-solid",
			Uri:      "admin/roles/index",
		},
		{
			BaseModel: model.BaseModel{
				ID: 4,
			},
			ParentId: 1,
			Order:    4,
			Name:     "权限",
			Icon:     "el-icon-lock",
			Uri:      "admin/permissions/index",
		},
		{
			BaseModel: model.BaseModel{
				ID: 5,
			},
			ParentId: 1,
			Order:    5,
			Name:     "菜单",
			Icon:     "el-icon-menu",
			Uri:      "admin/menus/index",
		},
		{
			BaseModel: model.BaseModel{
				ID: 6,
			},
			ParentId: 1,
			Order:    6,
			Name:     "日志",
			Icon:     "el-icon-notebook-2",
			Uri:      "admin/logs/index",
		},
		{
			BaseModel: model.BaseModel{
				ID: 7,
			},
			ParentId: 1,
			Order:    7,
			Name:     "文件",
			Icon:     "el-icon-files",
			Uri:      "admin/files/index",
		},
		{
			BaseModel: model.BaseModel{
				ID: 8,
			},
			ParentId: 1,
			Order:    8,
			Name:     "配置",
			Icon:     "el-icon-setting",
			Uri:      "admin/configs/index",
		},
	}

	return insertIgnoreOrUpdate(menus, false)
}
