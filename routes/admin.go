// Package routes 注册路由
package routes

import (
	"github.com/zero7cola/gin-admin-core/controllers"
	middlewares "github.com/zero7cola/gin-admin-core/middlerwares"

	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes 注册 admin 相关路由
func RegisterAdminRoutes(root *gin.RouterGroup) {
	var admin *gin.RouterGroup

	admin = root

	// 全局限流中间件：每小时限流。这里是所有 API （根据 IP）请求加起来。
	// 作为参考 Github API 每小时最多 60 个请求（根据 IP）。
	// 测试时，可以调高一点。
	//admin.Use(middlewares.LimitIP("500-H"))
	//admin.Use(middlewares.AuthAdminJWT())
	{

		ic := new(controllers.AdminIndexController)
		admin.GET("/index", ic.Index)
		admin.GET("/version", ic.Version)

		authGroup := admin.Group("/auth")
		// 登录
		lgc := new(controllers.AdminAuthController)
		authGroup.POST("/login", middlewares.GuestJWT(), lgc.Login)
		authGroup.POST("/refresh-token", lgc.RefreshToken)
		authGroup.GET("/current", lgc.Current)
		authGroup.POST("/logout", lgc.Logout)
		authGroup.GET("/captcha", lgc.ShowCaptcha)

		auc := new(controllers.AdminUserController)
		// 账号
		admin.GET("/users", auc.Index)
		admin.GET("/user/:id", auc.Get)
		admin.POST("/user", auc.Store)
		admin.PUT("/user/:id", auc.Update)
		admin.DELETE("/user/:id", auc.Delete)

		// 角色
		arc := new(controllers.AdminRoleController)
		admin.GET("/roles", arc.Index)
		admin.GET("/roles/all", arc.All)
		admin.GET("/role/:id", arc.Get)
		admin.POST("/role", arc.Store)
		admin.PUT("/role/:id", arc.Update)
		admin.DELETE("/role/:id", arc.Delete)

		// 菜单
		amc := new(controllers.AdminMenuController)
		admin.GET("/menus", amc.Index)
		admin.GET("/menus/all", amc.All)
		admin.GET("/menu/:id", amc.Get)
		admin.POST("/menu", amc.Store)
		admin.PUT("/menu/:id", amc.Update)
		admin.DELETE("/menu/:id", amc.Delete)

		// 权限
		apc := new(controllers.AdminPermissionController)
		admin.GET("/permissions", apc.Index)
		admin.GET("/permissions/all", apc.All)
		admin.GET("/permission/:id", apc.Get)
		admin.POST("/permission", apc.Store)
		admin.PUT("/permission/:id", apc.Update)
		admin.DELETE("/permission/:id", apc.Delete)

		// 配置
		cc := new(controllers.AdminConfigController)
		admin.GET("/configs", cc.Index)
		admin.GET("/configs/all", cc.All)
		admin.GET("/config/:id", cc.Get)
		admin.POST("/config", cc.Store)
		admin.PUT("/config/:id", cc.Update)
		admin.DELETE("/config/:id", cc.Delete)

		fc := new(controllers.AdminFileController)
		admin.POST("/upload", fc.Upload)
		admin.POST("/file", fc.Store)
		admin.GET("/files", fc.Index)
		admin.GET("/file/:id", fc.Get)
		admin.PUT("/file/:id", fc.Update)
		admin.DELETE("/file/:id", fc.Delete)

	}

}
