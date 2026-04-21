package adminPermission

import (
	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/core"
	"github.com/zero7cola/gin-admin-core/internal"
	"github.com/zero7cola/gin-admin-core/model"
	"github.com/zero7cola/gin-admin-core/pkg/paginator"
)

type AdminPermission struct {
	model.BaseModel
	Name       string `json:"name" gorm:"name"`
	Slug       string `json:"slug" gorm:"slug"`
	HttpMethod string `json:"http_method" gorm:"http_method"`
	HttpPath   string `json:"http_path" gorm:"http_path"`
	Order      uint64 `json:"order" gorm:"order"`
	ParentId   uint64 `json:"parent_id" gorm:"parent_id"`
	model.CommonTimestampsField
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (model *AdminPermission) Create() {
	core.Global.DB.Create(&model)
}

func (model *AdminPermission) Save() (rowsAffected int64) {
	result := core.Global.DB.Save(&model)
	return result.RowsAffected
}

func (model *AdminPermission) Delete() (rowsAffected int64) {
	result := core.Global.DB.Delete(&model)
	return result.RowsAffected
}

func All() (models []AdminPermission) {
	core.Global.DB.Find(&models)
	return
}

func Get(idstr string) (userModel AdminPermission) {
	core.Global.DB.Where("id", idstr).First(&userModel)
	return
}

// Paginate 分页内容
func Paginate(c *gin.Context, perPage int) (users []AdminPermission, paging paginator.Paging) {
	paging = paginator.Paginate(
		c,
		core.Global.DB.Model(AdminPermission{}),
		&users,
		internal.VADMINURL(internal.TableName(&AdminPermission{})),
		perPage,
	)
	return
}
