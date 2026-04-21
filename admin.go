package admin

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/core"
	middlewares "github.com/zero7cola/gin-admin-core/middlerwares"
	"github.com/zero7cola/gin-admin-core/routes"
)

var builtinModules = []core.Module{}

// 注册内置模块
func RegisterBuiltin(mods ...core.Module) {
	builtinModules = append(builtinModules, mods...)
}

func Register(r *gin.Engine, prefix string, modules ...core.Module) {

	// 注册全局中间件
	registerGlobalMiddleWare(r)

	if prefix == "" {
		prefix = "/admin"
	}

	root := r.Group(prefix)

	// 注册中间件
	//root.Use(middlewares.LimitIP("500-H"))
	root.Use(middlewares.AuthAdminJWT())
	root.Use(middlewares.OperationLog())

	// 1️⃣ 内置模块
	//for _, m := range builtinModules {
	//	registerModule(root, m)
	//}

	// 注册静态资源路由
	routes.RegisterStaticRoutes(r)

	routes.RegisterAdminRoutes(root)

	// 2️⃣ 业务模块
	for _, m := range modules {
		registerModule(root, m)
	}

	//  配置 404 路由
	setup404Handler(r)
}

func registerGlobalMiddleWare(router *gin.Engine) {
	router.Use(
		//gin.Logger(),
		middlewares.Logger(),
		middlewares.Recovery2(),
		//cors.Default(),
		//gin.Recovery(),
		//middlewares.ForceUA(),
		cors.New(cors.Config{
			AllowAllOrigins: true,
			//AllowOrigins:     []string{"http://localhost:4000"}, // 改成你的前端地址
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
		}),
	)
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

func setup404Handler(router *gin.Engine) {
	// 处理 404 请求
	router.NoRoute(func(c *gin.Context) {
		// 获取标头信息的 Accept 信息
		acceptString := c.Request.Header.Get("Accept")
		if strings.Contains(acceptString, "text/html") {
			// 如果是 HTML 的话
			c.String(http.StatusNotFound, "页面返回 404")
		} else {
			// 默认返回 JSON
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    404,
				"error_message": "路由未定义，请确认 url 和请求方法是否正确。",
			})
		}
	})
}
