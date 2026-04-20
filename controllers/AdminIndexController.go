package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/core"
)

type AdminIndexController struct {
	BaseAPIController
}

func (ic *AdminIndexController) Index(c *gin.Context) {

	fmt.Printf("%v \n", core.Global.Config)

	core.Global.Logger.Info("AdminIndexController Index")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}

func (ic *AdminIndexController) Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": core.Global.Config.App.Version,
	})
}
