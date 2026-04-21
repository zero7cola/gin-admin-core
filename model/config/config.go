package config

import (
	"github.com/zero7cola/gin-admin-core/core"
	"github.com/zero7cola/gin-admin-core/internal"
	"github.com/zero7cola/gin-admin-core/model"
	"github.com/zero7cola/gin-admin-core/pkg/paginator"

	"github.com/gin-gonic/gin"
)

type Config struct {
	model.BaseModel
	ConfigKey   string `gorm:"column:config_key" json:"config_key"`
	ConfigValue string `gorm:"column:config_value" json:"config_value"`
	ConfigLabel string `gorm:"column:config_label" json:"config_label"`
	Type        int    `gorm:"column:type" json:"type"`
	Options     string `gorm:"column:options" json:"options"`
	Describe    string `gorm:"column:describe" json:"describe"`
	IsCanFront  int    `gorm:"column:is_can_front" json:"is_can_front"`
	IsRequired  uint   `gorm:"column:is_required" json:"is_required"`
	Order       uint   `gorm:"column:order" json:"order"`
	GroupId     uint   `gorm:"column:group_id" json:"group_id"`
	State       uint   `gorm:"column:state" json:"state"`
	ShowType    string `gorm:"column:show_type" json:"show_type"`
	Placeholder string `gorm:"column:placeholder" json:"placeholder"`
	model.CommonTimestampsField
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (model *Config) Create() {
	core.Global.DB.Create(&model)
}

func (model *Config) Save() (rowsAffected int64) {
	result := core.Global.DB.Save(&model)
	return result.RowsAffected
}

func (model *Config) Delete() (rowsAffected int64) {
	result := core.Global.DB.Delete(&model)
	return result.RowsAffected
}

func All() (models []Config) {
	core.Global.DB.Find(&models)
	return
}

func AllShow() (models []Config) {
	core.Global.DB.Where("is_can_front = 1").Find(&models)
	return
}

func Get(idstr string) (model Config) {
	core.Global.DB.Where("id", idstr).First(&model)
	return
}

// Paginate 分页内容
func Paginate(c *gin.Context, perPage int) (users []Config, paging paginator.Paging) {

	db := core.Global.DB.Model(Config{})

	if c.Query("config_key") != "" {
		db = db.Where("config_key LIKE ?", "%"+c.Query("config_key")+"%")
	}

	if c.Query("config_label") != "" {
		db = db.Where("config_label LIKE ?", "%"+c.Query("config_label")+"%")
	}

	paging = paginator.Paginate(
		c,
		db,
		&users,
		internal.VADMINURL(internal.TableName(&Config{})),
		perPage,
	)
	return
}
