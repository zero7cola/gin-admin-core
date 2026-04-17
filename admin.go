package admin

import "github.com/gin-gonic/gin"

type Module func(rg *gin.RouterGroup)

// 内置模块（统一管理）
var builtinModules = []Module{
	// 在这里放内置模块
	// health.Register,
	// auth.Register,
}

// 提供方法注册内置模块（避免循环依赖）
func RegisterBuiltin(modules ...Module) {
	builtinModules = append(builtinModules, modules...)
}

// Register
func Register(r *gin.Engine, prefix string, modules ...Module) {
	if prefix == "" {
		prefix = "/admin"
	}

	root := r.Group(prefix)

	// 1️⃣ 先加载内置模块
	for _, m := range builtinModules {
		m(root)
	}

	// 2️⃣ 再加载业务模块
	for _, m := range modules {
		m(root)
	}
}
