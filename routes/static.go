// Package routes 注册路由
package routes

import (
	"github.com/zero7cola/gin-admin-core/config"

	"github.com/gin-gonic/gin"
)

// RegisterStaticRoutes 注册 static 相关路由
func RegisterStaticRoutes(r *gin.Engine) {

	// 本地文件
	r.Static("/"+config.GetString("storage.local.static"), "storage/files")

}
