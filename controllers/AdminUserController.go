package controllers

import (
	"fmt"
	"net/http"

	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/model/adminRole"
	"github.com/zero7cola/gin-admin-core/model/adminUser"
	"github.com/zero7cola/gin-admin-core/pkg/auth"
	"github.com/zero7cola/gin-admin-core/pkg/database"
	"github.com/zero7cola/gin-admin-core/pkg/helpers"
	"github.com/zero7cola/gin-admin-core/pkg/response"
	"github.com/zero7cola/gin-admin-core/requests"
)

type AdminUserController struct {
	BaseAPIController
}

func (uc *AdminUserController) Index(c *gin.Context) {

	data, pager := adminUser.Paginate(c, GetPerPage(c))

	response.Data(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

func (uc *AdminUserController) Get(c *gin.Context) {

	user := adminUser.Get(c.Param("id"))

	response.Data(c, user)
}

func (uc *AdminUserController) Store(c *gin.Context) {
	// 验证
	request := requests.AdminUserStoreRequest{}
	if ok := requests.ValidateFunc(c, &request, requests.VerityAdminUserStore); !ok {
		return
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
		}
	}()

	// 查询角色
	var roles []adminRole.AdminRole

	if len(request.RoleIDs) > 0 {
		if err := tx.Where("id IN ?", request.RoleIDs).Find(&roles).Error; err != nil {
			response.BadRequest(c, err, "角色查询失败")
			return
		}
	}

	model := adminUser.AdminUser{
		Username: request.Username,
		Password: request.Password,
		Name:     request.Name,
		Roles:    roles,
		AvatarId: helpers.Uint64Ptr(cast.ToUint64(request.AvatarId)),
	}

	if request.AvatarId > 0 {

	}

	if err := tx.Create(&model).Error; err != nil {
		tx.Rollback()
		response.BadRequest(c, err, "创建账号失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		response.BadRequest(c, err, "提交事务失败")
		return
	}

	response.Data(c, model)
}

func (uc *AdminUserController) Update(c *gin.Context) {
	userModel := adminUser.Get(c.Param("id"))
	if userModel.ID <= 0 {
		response.Fail(c, "没有找到")
		return
	}

	// 验证
	request := requests.AdminUserUpdateRequest{}
	request.ID = userModel.ID
	if ok := requests.ValidateFunc(c, &request, requests.VerityAdminUserUpdate); !ok {
		return
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
		}
	}()

	// 查询角色
	var roles []adminRole.AdminRole

	if len(request.RoleIDs) > 0 {
		if err := tx.Where("id IN ?", request.RoleIDs).Find(&roles).Error; err != nil {
			response.BadRequest(c, err, "角色查询失败")
			return
		}
	}

	// 替换关联角色
	if err := tx.Model(&userModel).Association("Roles").Replace(roles); err != nil {
		tx.Rollback()
		return
	}

	if !helpers.Empty(request.Username) {
		userModel.Username = request.Username
	}

	if !helpers.Empty(request.Password) {
		userModel.Password = request.Password
	}

	if !helpers.Empty(request.Name) {
		userModel.Name = request.Name
	}

	if request.AvatarId > 0 {
		userModel.AvatarId = helpers.Uint64Ptr(cast.ToUint64(request.AvatarId))
	}

	//fmt.Printf("%T", userModel)

	if err := tx.Save(&userModel).Error; err != nil {
		tx.Rollback()
		response.BadRequest(c, err, "创建账号失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		response.BadRequest(c, err, "提交事务失败")
		return
	}

	response.Data(c, userModel)
}

func (uc *AdminUserController) Delete(c *gin.Context) {

	fmt.Printf("current_user_id :%v", auth.CurrentAdminUser(c).ID)

	userModel := adminUser.Get(c.Param("id"))
	if userModel.ID <= 0 {
		response.Fail(c, "没有找到")
		return
	}

	if res := userModel.Delete(); res > 0 {
		response.Success(c)
		return
	}

	response.Fail(c, "删除失败")

}
