package internal

import (
	"sync"

	"github.com/zero7cola/gin-admin-core/setting"

	"time"
)

func IsLocal() bool {
	return setting.GlobalSetting.App.Env == "local"
}

func IsProduction() bool {
	return setting.GlobalSetting.App.Env == "production"
}

func IsTesting() bool {
	return setting.GlobalSetting.App.Env == "testing"
}

func IsDebug() bool {
	return setting.GlobalSetting.App.Debug == true
}

// TimenowInTimezone 获取当前时间，支持时区
func TimenowInTimezone() time.Time {
	chinaTimezone, _ := time.LoadLocation(setting.GlobalSetting.App.Timezone)
	return time.Now().In(chinaTimezone)
}

// URL 传参 path 拼接站点的 URL
func URL(path string) string {
	return setting.GlobalSetting.App.Url + "/" + path
}

// VADMINURL 拼接带 admin 标示 URL
func VADMINURL(path string) string {
	return URL("/admin/" + path)
}

// V1URL 拼接带 v1 标示 URL
func V1URL(path string) string {
	return URL("/v1/" + path)
}

type Router struct {
	mu          sync.RWMutex
	ignorePaths []string
}

// 初始值直接写死
var globalRouter = &Router{
	ignorePaths: []string{"/admin/auth/login", "/admin/auth/captcha", "/admin/upload", "/admin/version", "/admin/test"}, // 忽略的路径无需验证,
}

// 只暴露方法，不暴露结构体
func GetIgnorePaths() []string {
	globalRouter.mu.RLock()
	defer globalRouter.mu.RUnlock()

	res := make([]string, len(globalRouter.ignorePaths))
	copy(res, globalRouter.ignorePaths)
	return res
}

func AppendIgnorePaths(v []string) {
	globalRouter.mu.Lock()
	defer globalRouter.mu.Unlock()

	globalRouter.ignorePaths = append(globalRouter.ignorePaths, v...)
}
