package wechat

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/power"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/templateMessage/request"
	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/config"
	"net/http"
)

var OfficialAccountApp *officialAccount.OfficialAccount

func init() {
	GetApp()
}

func GetCallbackIP(ctx *gin.Context) {
	data, err := OfficialAccountApp.Base.GetCallbackIP(ctx)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func GetApp() {
	var w_err error
	OfficialAccountApp, w_err = officialAccount.NewOfficialAccount(&officialAccount.UserConfig{
		AppID:  config.GetString("offiaccount.appid"),     // 公众号、小程序的appid
		Secret: config.GetString("offiaccount.appsecret"), //
		Token:  config.GetString("offiaccount.token"),

		Log: officialAccount.Log{
			Level:  "debug",
			File:   "./wechat.log",
			Stdout: false, //  是否打印在终端
		},

		HttpDebug: true,
		Debug:     false,
	})

	if w_err != nil {
		panic(w_err)
	}
}

func SendWechatTemplateMessage(ctx *gin.Context, toUser string) {
	//OfficialAccountApp, err := officialAccount.NewOfficialAccount(&officialAccount.UserConfig{
	//	AppID:  config.GetString("wechat.offiaccount.appid"),     // 公众号、小程序的appid
	//	Secret: config.GetString("wechat.offiaccount.appsecret"), //
	//
	//	Log: officialAccount.Log{
	//		Level:  "debug",
	//		File:   "./wechat.log",
	//		Stdout: false, //  是否打印在终端
	//	},
	//
	//	HttpDebug: true,
	//	Debug:     false,
	//})

	//if err != nil {
	//	panic(err)
	//}

	OfficialAccountApp.TemplateMessage.Send(ctx, &request.RequestTemlateMessage{
		ToUser:     toUser,
		TemplateID: "templateID",
		URL:        "https://www.artisan-cloud.com/",
		Data: &power.HashMap{
			"first": &power.HashMap{
				"value": "恭喜你购买成功！",
				"color": "#173177",
			},
			"DateTime": &power.HashMap{
				"value": "2022-3-5 16:22",
				"color": "#173177",
			},
			"PayAmount": &power.HashMap{
				"value": "59.8元",
				"color": "#173177",
			},
			"Location": &power.HashMap{
				"value": "上海市长宁区",
				"color": "#173177",
			},
			"remark": &power.HashMap{
				"value": "欢迎再次购买！",
				"color": "#173177",
			},
		},
	})
}
