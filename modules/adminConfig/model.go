package adminConfig

import "github.com/zero7cola/gin-admin-core/modules/base"

type AdminConfig struct {
	base.BaseModel
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
	base.CommonTimestampsField
}
