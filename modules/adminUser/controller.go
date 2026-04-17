package adminUser

import (
	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/modules/base"
	"github.com/zero7cola/gin-admin-core/pkg/response"
)

type AdminUserController struct {
	base.BaseAPIController
}

func (uc *AdminUserController) Index(c *gin.Context) {

	response.Data(c, gin.H{
		"code": 200,
	})
}
