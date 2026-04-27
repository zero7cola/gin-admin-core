package adminPermission

import (
	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/internal"
	"github.com/zero7cola/gin-admin-core/model"
	"github.com/zero7cola/gin-admin-core/pkg/database"
	"github.com/zero7cola/gin-admin-core/pkg/paginator"
)

type AdminPermission struct {
	model.BaseModel
	Name       string `json:"name" gorm:"column:name;type:varchar(100)"`
	Slug       string `json:"slug" gorm:"column:slug;type:varchar(100)"`
	HttpMethod string `json:"http_method" gorm:"column:http_method;type:varchar(100)"`
	HttpPath   string `json:"http_path" gorm:"column:http_path;type:text"`
	Order      uint64 `json:"order" gorm:"column:order"`
	ParentId   uint64 `json:"parent_id" gorm:"column:parent_id"`
	model.CommonTimestampsField
}

func (model *AdminPermission) TableName() string {
	return "admin_permissions"
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (model *AdminPermission) Create() {
	database.DB.Create(&model)
}

func (model *AdminPermission) Save() (rowsAffected int64) {
	result := database.DB.Save(&model)
	return result.RowsAffected
}

func (model *AdminPermission) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&model)
	return result.RowsAffected
}

func All() (models []AdminPermission) {
	database.DB.Find(&models)
	return
}

func Get(idstr string) (userModel AdminPermission) {
	database.DB.Where("id", idstr).First(&userModel)
	return
}

// Paginate 分页内容
func Paginate(c *gin.Context, perPage int) (users []AdminPermission, paging paginator.Paging) {
	paging = paginator.Paginate(
		c,
		database.DB.Model(AdminPermission{}),
		&users,
		internal.VADMINURL(model.TableName(&AdminPermission{})),
		perPage,
	)
	return
}
