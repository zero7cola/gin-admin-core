package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/config"
	"net/http"
)

type AdminIndexController struct {
	BaseAPIController
}

func (ic *AdminIndexController) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}

func (ic *AdminIndexController) Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": config.GetString("version"),
	})
}
