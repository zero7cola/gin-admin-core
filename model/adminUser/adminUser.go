package adminUser

import (
	"github.com/zero7cola/gin-admin-core/internal"
	"github.com/zero7cola/gin-admin-core/model"
	"github.com/zero7cola/gin-admin-core/model/adminMenu"
	"github.com/zero7cola/gin-admin-core/model/adminPermission"
	"github.com/zero7cola/gin-admin-core/model/adminRole"
	"github.com/zero7cola/gin-admin-core/model/file"
	"github.com/zero7cola/gin-admin-core/pkg/database"
	"github.com/zero7cola/gin-admin-core/pkg/hash"
	"github.com/zero7cola/gin-admin-core/pkg/paginator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminUser struct {
	model.BaseModel
	Username      string                            `json:"username" gorm:"column:username;type:varchar(100)"`
	Password      string                            `json:"-" gorm:"column:password;type:varchar(255)"`
	Name          string                            `json:"name" gorm:"column:name;type:varchar(100)"`
	AvatarFile    *file.File                        `json:"avatar" gorm:"foreignKey:AvatarId;references:ID"`
	AvatarId      *uint64                           `json:"avatar_id" gorm:"column:avatar_id"`
	RememberToken string                            `json:"-" gorm:"column:remember_token;type:varchar(255)"`
	Roles         []adminRole.AdminRole             `json:"roles" gorm:"many2many:admin_role_users;foreignKey:ID;joinForeignKey:UserID;references:ID;joinReferences:RoleID"`
	Permissions   []adminPermission.AdminPermission `json:"permissions" gorm:"-"`
	Menus         []adminMenu.AdminMenu             `json:"menus" gorm:"-"`
	model.CommonTimestampsField
}

func (model *AdminUser) TableName() string {
	return "admin_users"
}

func GetUserPermissions(userID uint64) ([]adminPermission.AdminPermission, error) {
	var user AdminUser
	if err := database.DB.
		Preload("Roles.Permissions").
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		return nil, err
	}

	if user.IsSuperAdmin() {
		return adminPermission.All(), nil
	}

	permissionMap := make(map[uint64]adminPermission.AdminPermission)
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			permissionMap[perm.ID] = perm
		}
	}

	var permissions []adminPermission.AdminPermission
	for _, perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func GetUserMenus(userID uint64) ([]adminMenu.AdminMenu, error) {
	var user AdminUser
	if err := database.DB.
		Preload("Roles.Menus").
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		return nil, err
	}

	var menus []adminMenu.AdminMenu

	if user.IsSuperAdmin() {
		return adminMenu.All(), nil
	}

	menuMap := make(map[uint64]adminMenu.AdminMenu)
	for _, role := range user.Roles {
		for _, menu := range role.Menus {
			menuMap[menu.ID] = menu
		}
	}

	for _, perm := range menuMap {
		menus = append(menus, perm)
	}

	return menus, nil
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (model *AdminUser) Create() {
	database.DB.Create(&model)
}

func (model *AdminUser) Save() (rowsAffected int64) {
	result := database.DB.Save(&model)
	return result.RowsAffected
}

func (model *AdminUser) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&model)
	return result.RowsAffected
}

func (model *AdminUser) IsSuperAdmin() bool {
	return model.ID == 1 || model.ID == 5
}

func Get(idstr string) (model AdminUser) {
	database.DB.Where("id", idstr).Preload("Roles").Preload("AvatarFile").First(&model)
	return
}

// Paginate 分页内容
func Paginate(c *gin.Context, perPage int) (users []AdminUser, paging paginator.Paging) {
	paging = paginator.Paginate(
		c,
		database.DB.Model(AdminUser{}),
		&users,
		internal.VADMINURL(model.TableName(&AdminUser{})),
		perPage,
	)
	return
}

// ComparePassword 密码是否正确
func (model *AdminUser) ComparePassword(_password string) bool {
	return hash.BcryptCheckIn(_password, model.Password)
}

// GetByMulti 通过 手机号/Email/用户名 来获取用户
func GetByMulti(loginID string) (model AdminUser) {
	database.DB.
		Where("username = ?", loginID).
		First(&model)
	return
}

// BeforeSave GORM 的模型钩子，在创建和更新模型前调用
func (model *AdminUser) BeforeSave(tx *gorm.DB) (err error) {

	if !hash.BcryptIsHashed(model.Password) {
		model.Password = hash.BcryptHash(model.Password)
	}
	return
}
