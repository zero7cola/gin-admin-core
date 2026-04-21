// Package auth 授权相关逻辑
package auth

import (
	"errors"

	"github.com/zero7cola/gin-admin-core/core"
	"github.com/zero7cola/gin-admin-core/model/adminUser"
	"github.com/zero7cola/gin-admin-core/pkg/logger"

	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
)

// AttemptAdmin 尝试登录
func AttemptAdmin(email string, password string) (adminUser.AdminUser, error) {
	userModel := adminUser.GetByMulti(email)
	if userModel.ID == 0 {
		return adminUser.AdminUser{}, errors.New("账号不存在")
	}

	if !userModel.ComparePassword(password) {
		return adminUser.AdminUser{}, errors.New("密码错误")
	}

	return userModel, nil
}

// CurrentAdminUser 从 gin.context 中获取当前登录用户
func CurrentAdminUser(c *gin.Context) adminUser.AdminUser {
	userModel := adminUser.AdminUser{}
	userModel = adminUser.Get(cast.ToString(c.MustGet("current_admin_user_id")))

	if userModel.ID <= 0 {
		core.Global.Logger.Error(errors.New("无法获取用户").Error())
		//response.Fail(c, "没有找到")
		return adminUser.AdminUser{}
	}

	if c.GetInt("menus_on") > 0 || cast.ToInt(c.Query("menus_on")) > 0 {
		// 获取账号的显示菜单
		menus, errs := adminUser.GetUserMenus(userModel.ID)

		if errs != nil {
			logger.LogIf(errs)

			return userModel
		}
		userModel.Menus = menus
	}
	// db is now a *DB value
	return userModel
}

// CurrentUID 从 gin.context 中获取当前登录用户 ID

func CurrentAdminUID(c *gin.Context) string {
	return c.GetString("current_admin_user_id")
}
