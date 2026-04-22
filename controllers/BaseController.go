package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zero7cola/gin-admin-core/setting"
)

// BaseAPIController 基础控制器
type BaseAPIController struct {
}

func GetPerPage(c *gin.Context) int {
	key := setting.GlobalSetting.Paging.UrlQueryPerPage

	if len(key) == 0 {
		key = "per_page"
	}

	defaultPerPage := setting.GlobalSetting.Paging.PerPage

	if defaultPerPage <= 0 {
		defaultPerPage = 10
	}

	return cast.ToInt(c.DefaultQuery(key, cast.ToString(defaultPerPage)))
}
