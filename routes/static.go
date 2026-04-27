// Package routes 注册路由
package routes

import (
	"fmt"

	"github.com/zero7cola/gin-admin-core/setting"

	"github.com/gin-gonic/gin"
)

// RegisterStaticRoutes 注册 static 相关路由
func RegisterStaticRoutes(r *gin.Engine) {

	staticPath := "static"

	if setting.GlobalSetting.Storage.Local.StaticPrefix != "" {
		staticPath = setting.GlobalSetting.Storage.Local.StaticPrefix
	}

	staticRoute := fmt.Sprintf("/%s", staticPath)

	if HasRoute(r, staticRoute) {
		// 本地文件
		r.Static(staticRoute, "storage/files")
	}

}
