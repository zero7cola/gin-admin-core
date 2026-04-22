package controllers

import (
	configModel "github.com/zero7cola/gin-admin-core/model/config"
	"github.com/zero7cola/gin-admin-core/pkg/response"
	"github.com/zero7cola/gin-admin-core/requests"

	"github.com/gin-gonic/gin"
)

type AdminConfigController struct {
	BaseAPIController
}

func (uc *AdminConfigController) Index(c *gin.Context) {

	data, pager := configModel.Paginate(c, GetPerPage(c))

	response.Data(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

func (uc *AdminConfigController) All(c *gin.Context) {
	response.Data(c, configModel.All())
}

func (uc *AdminConfigController) Get(c *gin.Context) {
	response.Data(c, configModel.Get(c.Param("id")))
}

func (uc *AdminConfigController) Store(c *gin.Context) {
	// 验证
	request := requests.ConfigModelStoreRequest{}
	if ok := requests.ValidateFunc(c, &request, requests.VerityConfigModelStore); !ok {
		return
	}

	u := configModel.Config{
		ConfigLabel: request.ConfigLabel,
		ConfigKey:   request.ConfigKey,
		ConfigValue: request.ConfigValue,
		Options:     request.Options,
		Type:        request.Type,
		Describe:    request.Describe,
		IsCanFront:  request.IsCanFront,
		Order:       request.Order,
		GroupId:     request.GroupId,
		State:       request.State,
		ShowType:    request.ShowType,
		Placeholder: request.Placeholder,
		IsRequired:  request.IsRequired,
	}

	u.Create()

	response.Data(c, u)
}

func (uc *AdminConfigController) Update(c *gin.Context) {
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

func (uc *AdminConfigController) Delete(c *gin.Context) {
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
