package controllers

import (
	"fmt"
	"github.com/zero7cola/gin-admin-core/setting"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/core"
)

type AdminIndexController struct {
	BaseAPIController
}

func (ic *AdminIndexController) Index(c *gin.Context) {

	core.Global.Logger.Info("AdminIndexController Index")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}

func (ic *AdminIndexController) Version(c *gin.Context) {
	fmt.Printf("%v \n", setting.GlobalSetting)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": setting.GlobalSetting.App.Version,
	})
}
