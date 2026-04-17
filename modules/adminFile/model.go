package adminFile

import (
	"time"

	"github.com/zero7cola/gin-admin-core/modules/base"
)

type File struct {
	base.BaseModel
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
	base.CommonTimestampsField
}
