package internal

import (
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
