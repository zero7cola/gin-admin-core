package controllers

import (
	"fmt"
	"net/http"

	"github.com/zero7cola/gin-admin-core/pkg/logger"
	"github.com/zero7cola/gin-admin-core/setting"

	"github.com/gin-gonic/gin"
)

type AdminIndexController struct {
	BaseAPIController
}

func (ic *AdminIndexController) Index(c *gin.Context) {

	logger.Info("AdminIndexController Index")

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
