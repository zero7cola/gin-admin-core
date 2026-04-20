package app

import (
	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/core"
	"github.com/zero7cola/gin-admin-core/modules/adminUser"
)

var builtinModules = []core.Module{}

// 注册内置模块
func RegisterBuiltin(mods ...core.Module) {
	builtinModules = append(builtinModules, mods...)
}

func Register(r *gin.Engine, prefix string, modules ...core.Module) {
	if prefix == "" {
		prefix = "/admin"
	}

	root := r.Group(prefix)

	// 1️⃣ 内置模块
	for _, m := range builtinModules {
		registerModule(root, m)
	}

	// 2️⃣ 业务模块
	for _, m := range modules {
		registerModule(root, m)
	}
}

// 核心逻辑
func registerModule(root *gin.RouterGroup, m core.Module) {
	prefix := m.Prefix()

	// 自动补 /
	if prefix != "" && prefix[0] != '/' {
		prefix = "/" + prefix
	}

	group := root
	if prefix != "" {
		group = root.Group(prefix)
	}

	m.Register(group)
}

func init() {
	RegisterBuiltin(
		&adminUser.Module{},
	)
}
