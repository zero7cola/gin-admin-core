package adminRole

import (
	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/modules/base"

	"github.com/zero7cola/gin-admin-core/modules/adminMenu"
	"github.com/zero7cola/gin-admin-core/modules/adminPermission"
	"github.com/zero7cola/gin-admin-core/pkg/app"
	"github.com/zero7cola/gin-admin-core/pkg/database"
	"github.com/zero7cola/gin-admin-core/pkg/paginator"
)

type AdminRole struct {
	base.BaseModel
	Name        string                            `json:"name" gorm:"name"`
	Slug        string                            `json:"slug" gorm:"slug"`
	Permissions []adminPermission.AdminPermission `json:"permissions" gorm:"many2many:admin_role_permissions;foreignKey:ID;joinForeignKey:RoleID;references:ID;joinReferences:PermissionID"`
	Menus       []adminMenu.AdminMenu             `json:"menus" gorm:"many2many:admin_role_menus;foreignKey:ID;joinForeignKey:RoleID;references:ID;joinReferences:MenuID"`
	base.CommonTimestampsField
}
