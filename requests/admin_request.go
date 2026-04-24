package requests

import (
	"time"

	"github.com/zero7cola/gin-admin-core/model"
)

// 账号创建
type AdminUserStoreRequest struct {
	Username        string `json:"username" validate:"required,unique=admin_users;username"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	Name            string `json:"name" validate:"required"`
	RoleIDs         []uint `json:"role_ids"`
}

func VerityAdminUserStore(obj interface{}) map[string][]string {

	messages := map[string]map[string]string{
		"Username": {
			"required": "用户名为必填项，参数名称 username",
			"unique":   "用户名已经存在",
		},
		"Password": {
			"required": "密码为必填项，参数名称 password",
			"min":      "密码最少为 6 位",
		},
		"ConfirmPassword": {
			"required": "确认密码为必填项，参数名称 confirm_password",
			"eqfield":  "与 password 不一致",
		},
		"Name": {
			"required": "昵称为必填项，参数名称 name",
		},
	}

	errors := ValidateStruct(obj, messages)

	return errors
}

// 账号更新
type AdminUserUpdateRequest struct {
	model.BaseModel
	Username string `json:"username" validate:"omitempty,unique=admin_users;username;id;ID"`
	Password string `json:"password" validate:"omitempty,min=6"`
	Name     string `json:"name"`
	RoleIDs  []uint `json:"role_ids"`
}

func VerityAdminUserUpdate(obj interface{}) map[string][]string {

	messages := map[string]map[string]string{
		"Username": {
			"unique": "用户名已经存在",
		},
		"Password": {
			"min": "密码最少为 6 位",
		},
	}

	errors := ValidateStruct(obj, messages)

	return errors
}

// 角色创建
type AdminRoleStoreRequest struct {
	Name          string `json:"name" validate:"required"`
	Slug          string `json:"slug" validate:"required,unique=admin_roles;slug"`
	PermissionIDs []uint `json:"permission_ids"`
	MenuIDs       []uint `json:"menu_ids"`
}

func VerityAdminRoleStore(obj interface{}) map[string][]string {

	messages := map[string]map[string]string{
		"Name": {
			"required": "角色名为必填项，参数名称 name",
		},
		"Slug": {
			"required": "标记名 为必填项，参数名称 slug",
			"unique":   "标记名 已经存在",
		},
	}

	errors := ValidateStruct(obj, messages)

	return errors
}

// 角色更新
type AdminRoleUpdateRequest struct {
	model.BaseModel        // 包含 unique 验证id规则的需要添加
	Name            string `json:"name" validate:"required"`
	Slug            string `json:"slug" validate:"required,unique=admin_roles;slug;id;ID"`
	PermissionIDs   []uint `json:"permission_ids"`
	MenuIDs         []uint `json:"menu_ids"`
}

func VerityAdminRoleUpdate(obj interface{}) map[string][]string {

	messages := map[string]map[string]string{
		"Name": {
			"required": "角色名 为必填项，参数名称 name",
		},
		"Slug": {
			"required": "标记名 为必填项，参数名称 slug",
			"unique":   "标记名 已经存在",
		},
	}

	errors := ValidateStruct(obj, messages)

	return errors
}

// 菜单创建
type AdminMenuStoreRequest struct {
	Name     string `json:"name" validate:"required,unique=admin_menus;name"`
	Order    uint64 `json:"order" validate:"required,numeric"`
	Uri      string `json:"uri" validate:"required"`
	ParentId uint64 `json:"parent_id" validate:"omitempty,numeric"`
	Icon     string `json:"icon"`
}

func VerityAdminMenuStore(obj interface{}) map[string][]string {

	messages := map[string]map[string]string{
		"Name": {
			"required": "菜单名 为必填项，参数名称 name",
			"unique":   "菜单名 已经存在",
		},
		"Uri": {
			"required": "路径 为必填项，参数名称 uri",
		},
		"Order": {
			"required": "排序 为必填项，参数名称 order",
			"numeric":  "排序 为数字类型",
		},
		"ParentId": {
			"numeric": "父菜单id 为数字类型",
		},
	}

	errors := ValidateStruct(obj, messages)

	return errors
}

// 菜单更新
type AdminMenuUpdateRequest struct {
	model.BaseModel        // 包含 unique 验证id规则的需要添加
	Name            string `json:"name" validate:"unique=admin_menus;name;id;ID"`
	Order           uint64 `json:"order" validate:"numeric"`
	ParentId        uint64 `json:"parent_id" validate:"numeric"`
	Icon            string `json:"icon"`
	Uri             string `json:"uri"`
}

func VerityAdminMenuUpdate(obj interface{}) map[string][]string {

	messages := map[string]map[string]string{
		"Name": {
			"unique": "菜单名 已经存在",
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

// 权限创建
type AdminPermissionStoreRequest struct {
	Name       string `json:"name" validate:"required"`
	Slug       string `json:"slug" validate:"required,unique=admin_permissions;slug"`
	HttpMethod string `json:"http_method" validate:"required"`
	HttpPath   string `json:"http_path" validate:"required"`
	Order      uint64 `json:"order" validate:"required,numeric"`
	ParentId   uint64 `json:"parent_id" validate:"omitempty,numeric"`
}

func VerityAdminPermissionStore(obj interface{}) map[string][]string {

	messages := map[string]map[string]string{
		"Name": {
			"required": "权限名 为必填项，参数名称 name",
		},
		"Slug": {
			"required": "标记名 为必填项，参数名称 slug",
			"unique":   "标记名 已经存在",
		},
		"Order": {
			"required": "排序 为必填项，参数名称 order",
			"numeric":  "排序 为数字类型",
		},
		"ParentId": {
			"numeric": "父菜单id 为数字类型",
		},
	}

	errors := ValidateStruct(obj, messages)

	return errors
}

// 权限更新
type AdminPermissionUpdateRequest struct {
	model.BaseModel        // 包含 unique 验证id规则的需要添加
	Name            string `json:"name" validate:"required"`
	Slug            string `json:"slug" validate:"required,unique=admin_permissions;slug;id;ID"`
	HttpMethod      string `json:"http_method"`
	HttpPath        string `json:"http_path"`
	Order           uint64 `json:"order" validate:"required,numeric"`
	ParentId        uint64 `json:"parent_id" validate:"omitempty,numeric"`
}

func VerityAdminPermissionUpdate(obj interface{}) map[string][]string {

	messages := map[string]map[string]string{
		"Title": {
			"unique": "菜单名 已经存在",
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

// 登录
type AdminLoginRequest struct {
	Username      string `json:"username" validate:"required"`
	Password      string `json:"password" validate:"required"`
	CaptchaID     string `json:"captcha_id,omitempty" validate:"required"`
	CaptchaAnswer string `json:"captcha_answer,omitempty" validate:"required"`
}

func VerityAdminLogin(obj interface{}) map[string][]string {

	messages := map[string]map[string]string{
		"Username": {
			"required": "用户名为必填项，参数名称 username",
		},
		"Password": {
			"required": "密码为必填项，参数名称 password",
		},
		"CaptchaID": {
			"required": "图片验证码的 ID 为必填",
		},
		"CaptchaAnswer": {
			"required": "图片验证码答案必填",
		},
	}

	errors := ValidateStruct(obj, messages)
	// 手机验证码
	_data := obj.(*AdminLoginRequest)
	errors = ValidateCaptcha(_data.CaptchaID, _data.CaptchaAnswer, errors)

	return errors
}

// 菜单创建
type AdminFileStoreRequest struct {
	OriginName   string     `json:"origin_name"`
	Name         string     `json:"name"`
	Key          string     `json:"key"`
	GroupId      int        `json:"group_id"`
	Size         int64      `json:"size"`
	Storage      string     `json:"storage"`
	Path         string     `json:"path"`
	Type         int        `json:"type"`
	Ext          string     `json:"ext"`
	UserId       uint64     `json:"user_id"`
	Url          string     `json:"url"`
	ContentType  string     `json:"content_type"`
	ETag         string     `json:"e_tag"`
	Bucket       string     `json:"bucket"`
	LastModified time.Time  `json:"last_modified"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

func VerityAdminFileStore(obj interface{}) map[string][]string {

	messages := map[string]map[string]string{
		"Name": {
			"required": "菜单名 为必填项，参数名称 name",
			"unique":   "菜单名 已经存在",
		},
		"Uri": {
			"required": "路径 为必填项，参数名称 uri",
		},
		"Order": {
			"required": "排序 为必填项，参数名称 order",
			"numeric":  "排序 为数字类型",
		},
		"ParentId": {
			"numeric": "父菜单id 为数字类型",
		},
	}

	errors := ValidateStruct(obj, messages)

	return errors
}

// 文件更新
type AdminFileUpdateRequest struct {
	model.BaseModel            // 包含 unique 验证id规则的需要添加
	OriginName      string     `json:"origin_name"`
	Name            string     `json:"name"`
	Key             string     `json:"key"`
	GroupId         int        `json:"group_id"`
	Size            int64      `json:"size"`
	Storage         string     `json:"storage"`
	Path            string     `json:"path"`
	Type            int        `json:"type"`
	Ext             string     `json:"ext"`
	UserId          uint64     `json:"user_id"`
	Url             string     `json:"url"`
	ContentType     string     `json:"content_type"`
	ETag            string     `json:"e_tag"`
	Bucket          string     `json:"bucket"`
	LastModified    time.Time  `json:"last_modified"`
	DeletedAt       *time.Time `json:"deleted_at"`
}

func VerityAdminFileUpdate(obj interface{}) map[string][]string {

	messages := map[string]map[string]string{
		"Name": {
			"unique": "菜单名 已经存在",
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

type AdminUserProfileUpdateRequest struct {
	Name            string `json:"name" validate:"required"`
	Password        string `json:"password" validate:"min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
}

func VerityAdminUserProfileUpdate(obj interface{}) map[string][]string {

	messages := map[string]map[string]string{
		"Name": {
			"required": "昵称为必填项，参数名称 name",
		},
		"Password": {
			"min": "密码最少为 6 位",
		},
		"ConfirmPassword": {
			"eqfield": "与 password 不一致",
		},
	}

	errors := ValidateStruct(obj, messages)

	return errors
}
