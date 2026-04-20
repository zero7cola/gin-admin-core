package requests

import "github.com/zero7cola/gin-admin-core/model"

type ConfigModelStoreRequest struct {
	ConfigKey   string `json:"config_key" validate:"required,unique=configs;config_key"`
	ConfigValue string `json:"config_value" validate:"required"`
	Type        int    `json:"type"`
	Options     string `json:"options"`
	Describe    string `json:"describe"`
	IsCanFront  int    `json:"is_can_front"`
	IsRequired  uint   `json:"is_required"`
	Order       uint   `json:"order"`
	GroupId     uint   `json:"group_id"`
	State       uint   `json:"state"`
	ShowType    string `json:"show_type"`
	ConfigLabel string `json:"config_label" validate:"required"`
	Placeholder string `json:"placeholder"`
}

func VerityConfigModelStore(obj interface{}) map[string][]string {

	messages := map[string]map[string]string{
		"ConfigLabel": {
			"required": "Label 为必填项，参数名称 label",
		},
		"ConfigKey": {
			"required": "key 为必填项，参数名称 key",
			"unique":   "Key 已经存在",
		},
		"ConfigValue": {
			"required": "Value 为必填项，参数名称 Value",
		},
	}

	errors := ValidateStruct(obj, messages)

	return errors
}

type ConfigModelUpdateRequest struct {
	model.BaseModel        // 包含 unique 验证id规则的需要添加
	ConfigKey       string `json:"config_key" validate:"unique=configs;config_key;id;ID"`
	ConfigValue     string `json:"config_value" validate:"required"`
	Type            int    `json:"type"`
	Options         string `json:"options"`
	Describe        string `json:"describe"`
	IsCanFront      int    `json:"is_can_front"`
	IsRequired      uint   `json:"is_required"`
	Order           uint   `json:"order"`
	GroupId         uint   `json:"group_id"`
	State           uint   `json:"state"`
	ShowType        string `json:"show_type"`
	ConfigLabel     string `json:"config_label"`
	Placeholder     string `json:"placeholder"`
}

func VerityConfigModelUpdate(obj interface{}) map[string][]string {

	messages := map[string]map[string]string{
		"ConfigKey": {
			"unique": "ConfigKey 已经存在",
		},
		"Order": {
			"numeric": "排序 数字类型",
		},
		"ParentId": {
			"numeric": "父菜单id 为数字类型",
		},
	}

	errors := ValidateStruct(obj, messages)

	return errors
}
