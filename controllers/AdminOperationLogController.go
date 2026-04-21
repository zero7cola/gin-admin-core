package controllers

import (
	"github.com/spf13/cast"
	"github.com/zero7cola/gin-admin-core/model/adminOperationLog"
	configModel "github.com/zero7cola/gin-admin-core/model/config"
	"github.com/zero7cola/gin-admin-core/pkg/auth"
	"github.com/zero7cola/gin-admin-core/pkg/response"
	"github.com/zero7cola/gin-admin-core/requests"

	"github.com/gin-gonic/gin"
)

type AdminOperationController struct {
	BaseAPIController
}

func (uc *AdminOperationController) Index(c *gin.Context) {

	data, pager := adminOperationLog.Paginate(c, 5)

	response.Data(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

func (uc *AdminOperationController) All(c *gin.Context) {
	response.Data(c, adminOperationLog.All())
}

func (uc *AdminOperationController) Get(c *gin.Context) {
	response.Data(c, adminOperationLog.Get(c.Param("id")))
}

func (uc *AdminOperationController) Store(c *gin.Context) {

	u := adminOperationLog.AdminOperationLog{
		UserId: cast.ToUint64(auth.CurrentAdminUID(c)),
	}

	u.Create()

	response.Data(c, u)
}

func (uc *AdminOperationController) Update(c *gin.Context) {
	model := configModel.Get(c.Param("id"))
	if model.ID <= 0 {
		response.Fail(c, "没有找到")
		return
	}

	// 验证
	request := requests.ConfigModelUpdateRequest{}
	request.ID = model.ID
	if ok := requests.ValidateFunc(c, &request, requests.VerityConfigModelUpdate); !ok {
		return
	}

	model.ConfigLabel = request.ConfigLabel
	model.ConfigKey = request.ConfigKey
	model.ConfigValue = request.ConfigValue
	model.Options = request.Options
	model.Type = request.Type
	model.Describe = request.Describe
	model.IsCanFront = request.IsCanFront
	model.Order = request.Order
	model.GroupId = request.GroupId
	model.State = request.State
	model.ShowType = request.ShowType
	model.Placeholder = request.Placeholder
	model.IsRequired = request.IsRequired

	model.Save()

	response.Data(c, model)
}

func (uc *AdminOperationController) Delete(c *gin.Context) {
	model := configModel.Get(c.Param("id"))
	if model.ID <= 0 {
		response.Fail(c, "没有找到")
		return
	}

	if res := model.Delete(); res > 0 {
		response.Success(c)
		return
	}

	response.Fail(c, "删除失败")

}
