package adminOperationLog

import (
	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/internal"
	"github.com/zero7cola/gin-admin-core/model"
	"github.com/zero7cola/gin-admin-core/model/adminUser"
	"github.com/zero7cola/gin-admin-core/pkg/database"
	"github.com/zero7cola/gin-admin-core/pkg/paginator"
)

type AdminOperationLog struct {
	model.BaseModel
	AdminUser adminUser.AdminUser `json:"admin_user" gorm:"foreignKey:UserId;references:ID"`
	UserId    uint64              `json:"user_id" gorm:"user_id"`
	Path      string              `json:"path" gorm:"column:path,index"`
	Url       string              `json:"url" gorm:"url"`
	Method    string              `json:"method" gorm:"column:method,index"`
	Ip        string              `json:"ip" gorm:"column:ip,index"`
	Input     string              `json:"input" gorm:"type:text;column:input"`
	model.CommonTimestampsField
}

func (model *AdminOperationLog) TableName() string {
	return "admin_operation_logs"
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (model *AdminOperationLog) Create() {
	database.DB.Create(&model)
}

func (model *AdminOperationLog) Save() (rowsAffected int64) {
	result := database.DB.Save(&model)
	return result.RowsAffected
}

func (model *AdminOperationLog) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&model)
	return result.RowsAffected
}

func All() (models []AdminOperationLog) {
	database.DB.Find(&models)
	return
}

func Get(idstr string) (model AdminOperationLog) {
	database.DB.Where("id", idstr).Preload("AdminUser").First(&model)
	return
}

// Paginate 分页内容
func Paginate(c *gin.Context, perPage int) (data []AdminOperationLog, paging paginator.Paging) {
	db := database.DB.Model(AdminOperationLog{})

	if c.Query("path") != "" {
		db = db.Where("path LIKE ?", c.Query("path")+"%")
	}

	if c.Query("ip") != "" {
		db = db.Where("ip LIKE ?", c.Query("ip")+"%")
	}

	paging = paginator.Paginate(
		c,
		db,
		&data,
		internal.VADMINURL(model.TableName(&AdminOperationLog{})),
		perPage,
	)
	return
}
