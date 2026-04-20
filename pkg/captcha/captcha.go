// Package captcha 处理图片验证码逻辑
package captcha

import (
	"sync"

	"github.com/mojocn/base64Captcha"
	"github.com/zero7cola/gin-admin-core/core"
	"github.com/zero7cola/gin-admin-core/setting"
)

type Captcha struct {
	Base64Captcha *base64Captcha.Captcha
}

// once 确保 internalCaptcha 对象只初始化一次
var once sync.Once

// internalCaptcha 内部使用的 Captcha 对象
var internalCaptcha *Captcha

// NewCaptcha 单例模式获取
func NewCaptcha() *Captcha {
	once.Do(func() {
		// 初始化 Captcha 对象
		internalCaptcha = &Captcha{}

		// 使用全局 Redis 对象，并配置存储 Key 的前缀
		store := RedisStore{
			RedisClient: core.Global.Redis,
			KeyPrefix:   setting.GlobalSetting.App.Name + ":captcha:",
		}

		// 配置 base64Captcha 驱动信息
		driver := base64Captcha.NewDriverDigit(
			setting.GlobalSetting.Captcha.Height,
			setting.GlobalSetting.Captcha.Width,
			setting.GlobalSetting.Captcha.Length,
			setting.GlobalSetting.Captcha.Maxskew,
			setting.GlobalSetting.Captcha.Dotcount,
		)

		// 实例化 base64Captcha 并赋值给内部使用的 internalCaptcha 对象
		internalCaptcha.Base64Captcha = base64Captcha.NewCaptcha(driver, &store)
	})

	return internalCaptcha
}

// GenerateCaptcha 生成图片验证码
func (c *Captcha) GenerateCaptcha() (id string, b64s, answer string, err error) {
	return c.Base64Captcha.Generate()
}

// VerifyCaptcha 验证验证码是否正确
func (c *Captcha) VerifyCaptcha(id string, answer string) (match bool) {
	// 第三个参数是验证后是否删除，我们选择 false
	// 这样方便用户多次提交，防止表单提交错误需要多次输入图片验证码
	return c.Base64Captcha.Verify(id, answer, false)
}
