package adminRole

import (
	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/internal"
	"github.com/zero7cola/gin-admin-core/model"
	"github.com/zero7cola/gin-admin-core/model/adminMenu"
	"github.com/zero7cola/gin-admin-core/model/adminPermission"
	"github.com/zero7cola/gin-admin-core/pkg/database"
	"github.com/zero7cola/gin-admin-core/pkg/paginator"
)

type AdminRole struct {
	model.BaseModel
	Name        string                            `json:"name" gorm:"name"`
	Slug        string                            `json:"slug" gorm:"slug"`
	Permissions []adminPermission.AdminPermission `json:"permissions" gorm:"many2many:admin_role_permissions;foreignKey:ID;joinForeignKey:RoleID;references:ID;joinReferences:PermissionID"`
	Menus       []adminMenu.AdminMenu             `json:"menus" gorm:"many2many:admin_role_menus;foreignKey:ID;joinForeignKey:RoleID;references:ID;joinReferences:MenuID"`
	model.CommonTimestampsField
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (userModel *AdminRole) Create() {
	database.DB.Create(&userModel)
}

func (userModel *AdminRole) Save() (rowsAffected int64) {
	result := database.DB.Save(&userModel)
	return result.RowsAffected
}

func (userModel *AdminRole) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&userModel)
	return result.RowsAffected
}

func All() (models []AdminRole) {
	database.DB.Find(&models)
	return
}

func Get(idstr string) (userModel AdminRole) {
	database.DB.Where("id", idstr).Preload("Menus").Preload("Permissions").First(&userModel)
	return
}

// Paginate 分页内容
func Paginate(c *gin.Context, perPage int) (users []AdminRole, paging paginator.Paging) {
	paging = paginator.Paginate(
		c,
		database.DB.Model(AdminRole{}),
		&users,
		internal.VADMINURL(model.TableName(&AdminRole{})),
		perPage,
	)
	return
}
