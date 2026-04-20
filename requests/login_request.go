package requests

type LoginRequest struct {
	Username string `json:"username" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

//
//func VerityLogin(obj interface{}) map[string][]string {
//
//	messages := map[string]map[string]string{
//		"Username": {
//			"email":    "用户名必须为正常邮箱格式",
//			"required": "用户名为必填项，参数名称 username",
//		},
//		"Password": {
//			"required": "密码为必填项，参数名称 password",
//			"min=6":    "密码最少为 6 位",
//		},
//	}
//
//	errors := ValidateStructAll(obj, messages)
//
//	return errors
//}
