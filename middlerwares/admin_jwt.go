package middlewares

import (
	"fmt"
	"strings"

	"github.com/zero7cola/gin-admin-core/model/adminUser"
	"github.com/zero7cola/gin-admin-core/pkg/helpers"
	"github.com/zero7cola/gin-admin-core/pkg/jwt"
	"github.com/zero7cola/gin-admin-core/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthAdminJWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		/**
		c.FullPath()   // "/user/:id"   ✅ 路由模板
		c.Request.URL.Path  // "/user/123"  ✅ 实际路径
		c.Param("id")  // "123"         ✅ 路径参数
		c.Query("name") // "abc"        ✅ query参数
		*/
		path := c.FullPath()
		ignorePaths := []string{"/admin/auth/login", "/admin/auth/captcha", "/admin/upload", "/admin/version", "/admin/test"} // 忽略的路径无需验证
		ignorePermissionPaths := []string{"/admin/auth/logout", "/admin/auth/refresh-token", "/admin/roles/all", "/admin/permissions/all", "/admin/menus/all", "/admin/auth/current"}

		//if !app.IsProduction() {
		//	fmt.Printf("full :%s path :%s \n", c.FullPath(), c.Request.URL.Path)
		//}

		//
		if !helpers.StringContains(ignorePaths, path) {
			// 从标头 Authorization:Bearer xxxxx 中获取信息，并验证 JWT 的准确性
			claims, err := jwt.NewJWT().ParserToken(c)

			// JWT 解析失败，有错误发生
			if err != nil {
				response.AuthFail(c, fmt.Sprintf("请查看相关的接口认证文档 path：%s", path))
				return
			}

			// 判断是否是 admin token
			if claims.Type != jwt.ADMIN_TOKEN_TYPE {
				response.Abort403(c, "非法操作!!!")
				return
			}

			// JWT 解析成功，设置用户信息
			userModel := adminUser.Get(claims.UserID)
			if userModel.ID == 0 {
				response.AuthFail(c, "找不到对应用户，用户可能已删除")
				return
			}

			// 验证权限
			// 3. 比对请求的 Method + 路由模板
			isPass := false

			// 有些全局的操作也不做验证权限
			if helpers.StringContains(ignorePermissionPaths, path) {
				isPass = true
			} else {
				// 超级管理员
				if userModel.IsSuperAdmin() {
					isPass = true
				} else {
					perms, err := adminUser.GetUserPermissions(userModel.ID)
					if err != nil {
						response.BadRequest(c, err, "权限加载失败")
						return
					}

					reqMethod := c.Request.Method // GET, POST, PUT ...
					for _, perm := range perms {
						if perm.HttpMethod != "any" && strings.ToUpper(perm.HttpMethod) != strings.ToUpper(reqMethod) {
							continue
						}
						if helpers.IsPathAllowed(path, perm.HttpPath) {
							// 放行
							isPass = true
							break
						}
					}
				}
			}

			if isPass {
				// 将用户信息存入 gin.context 里，后续 auth 包将从这里拿到当前用户数据
				c.Set("current_admin_user_id", userModel.GetStringID())
				//c.Set("current_user_name", userModel.Name)
				//c.Set("current_user", userModel)

				c.Next()
				return
			}

			// 4. 拒绝访问
			response.Abort403(c, "无权访问")
			return
		}

		c.Next()
	}
}
