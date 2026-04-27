package config

import (
	"github.com/zero7cola/gin-admin-core/internal"
	"github.com/zero7cola/gin-admin-core/model"
	"github.com/zero7cola/gin-admin-core/pkg/database"
	"github.com/zero7cola/gin-admin-core/pkg/paginator"

	"github.com/gin-gonic/gin"
)

type Config struct {
	model.BaseModel
	ConfigKey   string `gorm:"column:config_key;index;type:varchar(100)" json:"config_key"`
	ConfigValue string `gorm:"column:config_value" json:"config_value"`
	ConfigLabel string `gorm:"column:config_label;type:varchar(100);index" json:"config_label"`
	Type        int    `gorm:"column:type" json:"type"`
	Options     string `gorm:"column:options;type:text" json:"options"`
	Describe    string `gorm:"column:describe" json:"describe"`
	IsCanFront  int    `gorm:"column:is_can_front" json:"is_can_front"`
	IsRequired  uint   `gorm:"column:is_required" json:"is_required"`
	Order       uint   `gorm:"column:order" json:"order"`
	GroupId     uint   `gorm:"column:group_id" json:"group_id"`
	State       uint   `gorm:"column:state" json:"state"`
	ShowType    string `gorm:"column:show_type;type:varchar(255)" json:"show_type"`
	Placeholder string `gorm:"column:placeholder;type:varchar(255)" json:"placeholder"`
	model.CommonTimestampsField
}

func (model *Config) TableName() string {
	return "configs"
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (model *Config) Create() {
	database.DB.Create(&model)
}

func (model *Config) Save() (rowsAffected int64) {
	result := database.DB.Save(&model)
	return result.RowsAffected
}

func (model *Config) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&model)
	return result.RowsAffected
}

func All() (models []Config) {
	database.DB.Find(&models)
	return
}

func AllShow() (models []Config) {
	database.DB.Where("is_can_front = 1").Find(&models)
	return
}

func Get(idstr string) (model Config) {
	database.DB.Where("id", idstr).First(&model)
	return
}

// Paginate 分页内容
func Paginate(c *gin.Context, perPage int) (users []Config, paging paginator.Paging) {

	db := database.DB.Model(Config{})

	if c.Query("config_key") != "" {
		db = db.Where("config_key LIKE ?", c.Query("config_key")+"%")
	}

	if c.Query("config_label") != "" {
		db = db.Where("config_label LIKE ?", c.Query("config_label")+"%")
	}

	paging = paginator.Paginate(
		c,
		db,
		&users,
		internal.VADMINURL(model.TableName(&Config{})),
		perPage,
	)
	return
}
