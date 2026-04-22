package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/model/adminMenu"
	"github.com/zero7cola/gin-admin-core/model/adminPermission"
	"github.com/zero7cola/gin-admin-core/model/adminRole"
	"github.com/zero7cola/gin-admin-core/pkg/database"
	"github.com/zero7cola/gin-admin-core/pkg/response"
	"github.com/zero7cola/gin-admin-core/requests"
)

type AdminRoleController struct {
	BaseAPIController
}

func (uc *AdminRoleController) Index(c *gin.Context) {

	data, pager := adminRole.Paginate(c, GetPerPage(c))

	response.Data(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

func (uc *AdminRoleController) All(c *gin.Context) {

	roles := adminRole.All()

	response.Data(c, roles)
}

func (uc *AdminRoleController) Get(c *gin.Context) {

	user := adminRole.Get(c.Param("id"))

	response.Data(c, user)
}

func (uc *AdminRoleController) Store(c *gin.Context) {
	// 验证
	request := requests.AdminRoleStoreRequest{}
	if ok := requests.ValidateFunc(c, &request, requests.VerityAdminRoleStore); !ok {
		return
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
		}
	}()

	// 查询菜单和权限
	var permissions []adminPermission.AdminPermission
	var menus []adminMenu.AdminMenu

	if len(request.PermissionIDs) > 0 {
		if err := tx.Where("id IN ?", request.PermissionIDs).Find(&permissions).Error; err != nil {
			response.BadRequest(c, err, "权限查询失败")
			return
		}
	}

	if len(request.MenuIDs) > 0 {
		if err := tx.Where("id IN ?", request.MenuIDs).Find(&menus).Error; err != nil {
			response.BadRequest(c, err, "菜单查询失败")
			return
		}
	}

	role := adminRole.AdminRole{
		Name:        request.Name,
		Slug:        request.Slug,
		Permissions: permissions,
		Menus:       menus,
	}

	if err := tx.Create(&role).Error; err != nil {
		tx.Rollback()
		response.BadRequest(c, err, "创建角色失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		response.BadRequest(c, err, "提交事务失败")
		return
	}

	response.Data(c, role)
}

func (uc *AdminRoleController) Update(c *gin.Context) {
	userModel := adminRole.Get(c.Param("id"))
	if userModel.ID <= 0 {
		response.Fail(c, "没有找到")
		return
	}

	// 验证
	request := requests.AdminRoleUpdateRequest{}
	request.ID = userModel.ID
	if ok := requests.ValidateFunc(c, &request, requests.VerityAdminRoleUpdate); !ok {
		return
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
		}
	}()

	// 查询菜单和权限
	var permissions []adminPermission.AdminPermission
	var menus []adminMenu.AdminMenu

	if len(request.PermissionIDs) > 0 {
		if err := tx.Where("id IN ?", request.PermissionIDs).Find(&permissions).Error; err != nil {
			response.BadRequest(c, err, "权限查询失败")
			return
		}
	}

	if len(request.MenuIDs) > 0 {
		if err := tx.Where("id IN ?", request.MenuIDs).Find(&menus).Error; err != nil {
			response.BadRequest(c, err, "菜单查询失败")
			return
		}
	}

	// 替换关联权限和菜单
	if err := tx.Model(&userModel).Association("Permissions").Replace(permissions); err != nil {
		tx.Rollback()
		return
	}
	if err := tx.Model(&userModel).Association("Menus").Replace(menus); err != nil {
		tx.Rollback()
		return
	}

	//role := adminRole.AdminRole{
	//	Name:        request.Name,
	//	Slug:        request.Slug,
	//	Permissions: permissions,
	//	Menus:       menus,
	//}

	//fmt.Printf("%v \n", request)

	userModel.Name = request.Name
	userModel.Slug = request.Slug
	userModel.Permissions = permissions
	userModel.Menus = menus

	//fmt.Printf("%T", userModel)

	if err := tx.Save(&userModel).Error; err != nil {
		tx.Rollback()
		response.BadRequest(c, err, "创建角色失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		response.BadRequest(c, err, "提交事务失败")
		return
	}

	response.Data(c, userModel)
}

func (uc *AdminRoleController) Delete(c *gin.Context) {
	userModel := adminRole.Get(c.Param("id"))
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
