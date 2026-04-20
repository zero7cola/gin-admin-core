package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zero7cola/gin-admin-core/pkg/captcha"
	"github.com/zero7cola/gin-admin-core/pkg/database"
	"github.com/zero7cola/gin-admin-core/pkg/response"
	"strings"
)

type ValidatorFunc func(interface{}) map[string][]string

func ValidateFunc(c *gin.Context, obj interface{}, handler ValidatorFunc) bool {

	// 1. 解析请求，支持 JSON 数据、表单请求和 URL Query
	if err := c.ShouldBind(obj); err != nil {
		response.BadRequest(c, err, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		return false
	}

	c.Set("validate_data", obj)

	// 2. 表单验证
	errs := handler(obj)

	// 3. 判断验证是否通过
	if len(errs) > 0 {
		response.ValidationFields(c, errs)
		return false
	}

	return true
}

// 通用验证函数
func ValidateStruct(data interface{}, messages map[string]map[string]string) map[string][]string {
	validate := validator.New()

	// 注册 unique 自定义验证器
	_ = validate.RegisterValidation("unique", func(fl validator.FieldLevel) bool {
		param := fl.Param() // admin_users;username;id;ID 或 admin_users;username
		parts := strings.Split(param, ";")

		// 基础参数
		tableName := parts[0]
		columnName := parts[1]

		var query = database.DB.Table(tableName).Where(columnName+" = ?", fl.Field().Interface())

		if len(parts) == 4 {
			// 更新时：排除自身 ID
			idColumn := parts[2]
			structIDField := parts[3]

			structVal := fl.Parent()
			idField := structVal.FieldByName(structIDField)
			if !idField.IsValid() {
				return false
			}

			query = query.Where(idColumn+" != ?", idField.Interface())
		}

		var count int64
		if err := query.Count(&count).Error; err != nil {
			return false
		}
		return count == 0
	})

	errs := make(map[string][]string)
	// 执行验证
	err := validate.Struct(data)
	if err == nil {
		return errs
	}

	// 转换为 map[string][]string
	for _, e := range err.(validator.ValidationErrors) {
		field := e.StructField()
		tag := e.Tag()

		// 获取自定义提示
		msg := tag
		if m, ok := messages[field]; ok {
			if customMsg, exists := m[tag]; exists {
				msg = customMsg
			}
		}

		errs[field] = append(errs[field], msg)
	}

	return errs
}

func ValidateCaptcha(captchaID, captchaAnswer string, errs map[string][]string) map[string][]string {
	if ok := captcha.NewCaptcha().VerifyCaptcha(captchaID, captchaAnswer); !ok {
		errs["captcha_answer"] = append(errs["captcha_answer"], "图片验证码错误")
	}
	return errs
}
