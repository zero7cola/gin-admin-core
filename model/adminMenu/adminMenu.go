package adminMenu

import (
	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/core"
	"github.com/zero7cola/gin-admin-core/internal"
	"github.com/zero7cola/gin-admin-core/model"
	"github.com/zero7cola/gin-admin-core/pkg/paginator"
)

type AdminMenu struct {
	model.BaseModel
	ParentId uint64 `json:"parent_id" gorm:"parent_id"`
	Order    uint64 `json:"order" gorm:"order"`
	Name     string `json:"name" gorm:"name"`
	Icon     string `json:"icon" gorm:"icon"`
	Uri      string `json:"uri" gorm:"uri"`
	model.CommonTimestampsField
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (model *AdminMenu) Create() {
	core.Global.DB.Create(&model)
}

func (model *AdminMenu) Save() (rowsAffected int64) {
	result := core.Global.DB.Save(&model)
	return result.RowsAffected
}

func (model *AdminMenu) Delete() (rowsAffected int64) {
	result := core.Global.DB.Delete(&model)
	return result.RowsAffected
}

func All() (models []AdminMenu) {
	core.Global.DB.Find(&models)
	return
}

func Get(idstr string) (userModel AdminMenu) {
	core.Global.DB.Where("id", idstr).First(&userModel)
	return
}

// Paginate 分页内容
func Paginate(c *gin.Context, perPage int) (users []AdminMenu, paging paginator.Paging) {
	paging = paginator.Paginate(
		c,
		core.Global.DB.Model(AdminMenu{}),
		&users,
		internal.VADMINURL(internal.TableName(&AdminMenu{})),
		perPage,
	)
	return
}
