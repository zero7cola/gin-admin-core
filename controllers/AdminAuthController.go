package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/internal"
	"github.com/zero7cola/gin-admin-core/model/adminUser"
	"github.com/zero7cola/gin-admin-core/pkg/auth"
	"github.com/zero7cola/gin-admin-core/pkg/captcha"
	"github.com/zero7cola/gin-admin-core/pkg/helpers"
	"github.com/zero7cola/gin-admin-core/pkg/jwt"
	"github.com/zero7cola/gin-admin-core/pkg/logger"
	"github.com/zero7cola/gin-admin-core/pkg/response"
	"github.com/zero7cola/gin-admin-core/requests"
)

type AdminAuthController struct {
	BaseAPIController
}

func (ac *AdminAuthController) Login(c *gin.Context) {
	// 验证
	request := requests.AdminLoginRequest{}
	if ok := requests.ValidateFunc(c, &request, requests.VerityAdminLogin); !ok {
		return
	}

	userModel, err := auth.AttemptAdmin(request.Username, request.Password)

	if err != nil {
		// 失败，显示错误提示
		response.Fail(c, "账号不存在或密码错误")
		return
	} else {
		token := jwt.NewJWT().IssueAdminToken(userModel.GetStringID(), userModel.Username)

		// 获取账号的显示菜单
		menus, errs := adminUser.GetUserMenus(userModel.ID)

		if errs != nil {
			response.Error(c, err, "获取玩家显示菜单错误")
			return
		}

		userModel.Menus = menus

		c.Set("current_admin_user_id", userModel.GetStringID())

		response.Data(c, gin.H{
			"token": token,
			"user":  userModel,
		})
	}
}

func (ac *AdminAuthController) Logout(c *gin.Context) {

	response.Success(c)
}

func (ac *AdminAuthController) Current(c *gin.Context) {

	user := auth.CurrentAdminUser(c)

	response.Data(c, user)
}

// RefreshToken 刷新 Access Token
func (ac *AdminAuthController) RefreshToken(c *gin.Context) {

	token, err := jwt.NewJWT().RefreshToken(c)

	if err != nil {
		response.Error(c, err, "令牌刷新失败")
	} else {
		response.Data(c, gin.H{
			"token": token,
		})
	}
}

// ShowCaptcha 显示图片验证码
func (ac *AdminAuthController) ShowCaptcha(c *gin.Context) {
	// 生成验证码
	id, b64s, answer, err := captcha.NewCaptcha().GenerateCaptcha()

	if internal.IsDebug() {
		fmt.Printf("获取验证码 id:%s answer:%s\n", id, answer)
	}

	// 记录错误日志，因为验证码是用户的入口，出错时应该记 error 等级的日志
	logger.LogIf(err)
	// 返回给用户
	response.Data(c, gin.H{
		"captcha_id":    id,
		"captcha_image": b64s,
	})
}

func (ac *AdminAuthController) UpdateProfile(c *gin.Context) {

	user := auth.CurrentAdminUser(c)

	// 验证
	request := requests.AdminUserProfileUpdateRequest{}
	if ok := requests.ValidateFunc(c, &request, requests.VerityAdminUserProfileUpdate); !ok {
		return
	}

	if !helpers.Empty(request.Name) {
		user.Name = request.Name
	}

	if !helpers.Empty(request.Password) {
		if request.Password == request.ConfirmPassword {
			user.Password = request.Password
		} else {
			response.Fail(c, "两次输入的密码不一致")
		}
	}

	user.Save()

	response.Data(c, user)
}
