package file

import (
	"strings"
	"time"

	"github.com/zero7cola/gin-admin-core/config"
	"github.com/zero7cola/gin-admin-core/model"
	"github.com/zero7cola/gin-admin-core/pkg/database"
	"github.com/zero7cola/gin-admin-core/pkg/helpers"
	"github.com/zero7cola/gin-admin-core/pkg/paginator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type File struct {
	model.BaseModel
	OriginName   string     `gorm:"column:origin_name" json:"origin_name"`
	Name         string     `gorm:"column:name" json:"name"`
	Key          string     `gorm:"column:key" json:"key"`
	GroupId      int        `gorm:"column:group_id;index" json:"group_id"`
	Size         int64      `gorm:"column:size" json:"size"`
	Storage      string     `gorm:"column:storage" json:"storage"`
	Path         string     `gorm:"column:path" json:"-"`
	Type         int        `gorm:"column:type" json:"type"`
	Ext          string     `gorm:"column:ext" json:"ext"`
	UserId       uint64     `gorm:"column:user_id" json:"-"`
	Url          string     `gorm:"column:url" json:"url"`
	ContentType  string     `gorm:"column:content_type" json:"content_type"`
	ETag         string     `gorm:"column:e_tag" json:"e_tag"`
	Bucket       string     `gorm:"column:bucket" json:"bucket"`
	LastModified time.Time  `gorm:"column:last_modified" json:"last_modified"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	model.CommonTimestampsField
}

// 查询后
func (model *File) AfterFind(tx *gorm.DB) (err error) {
	model.Url = model.GetFileFullUrl()
	return
}

// 获取文件完整访问路径
func (model *File) GetFileFullUrl() string {
	url := model.Url
	if helpers.Empty(url) {
		url = config.GetString("storage." + model.Storage + ".domain")
		if model.Storage == "local" {
			path := strings.ReplaceAll(model.Path, config.GetString("storage.local.path"), config.GetString("storage.local.static"))
			url = url + "/" + path
		} else {
			url = url + "/" + model.Path
		}
	}
	return url
}

func (model *File) GetFileFullPath() string {
	path := model.Path
	if model.Storage == "local" {
		path = config.GetString("storage.local.path") + "/" + path
	} else {
		path = model.Path
	}
	return path
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (model *File) Create() {
	database.DB.Create(&model)
}

func (model *File) Save() (rowsAffected int64) {
	result := database.DB.Save(&model)
	return result.RowsAffected
}

func (model *File) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&model)
	return result.RowsAffected
}

func All() (models []File) {
	database.DB.Find(&models)
	return
}

func Get(idstr string) (model File) {
	database.DB.Where("id", idstr).First(&model)
	return
}

// Paginate 分页内容
func Paginate(c *gin.Context, perPage int) (users []File, paging paginator.Paging) {
	db := database.DB.Model(File{})
	if c.Query("storage") != "" {
		db = db.Where("storage = ?", c.Query("storage"))
	}

	paging = paginator.Paginate(
		c,
		db,
		&users,
		helpers.VADMINURL(database.TableName(&File{})),
		perPage,
	)
	return
}
