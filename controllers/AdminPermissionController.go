package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/model/adminPermission"
	"github.com/zero7cola/gin-admin-core/pkg/response"
	"github.com/zero7cola/gin-admin-core/requests"
	"strings"
)

type AdminPermissionController struct {
	BaseAPIController
}

func (uc *AdminPermissionController) Index(c *gin.Context) {

	data, pager := adminPermission.Paginate(c, 10)

	response.Data(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

func (uc *AdminPermissionController) All(c *gin.Context) {

	models := adminPermission.All()

	response.Data(c, models)
}

func (uc *AdminPermissionController) Get(c *gin.Context) {

	user := adminPermission.Get(c.Param("id"))

	response.Data(c, user)
}

func (uc *AdminPermissionController) Store(c *gin.Context) {
	// 验证
	request := requests.AdminPermissionStoreRequest{}
	if ok := requests.ValidateFunc(c, &request, requests.VerityAdminPermissionStore); !ok {
		return
	}

	u := adminPermission.AdminPermission{
		Name:       request.Name,
		Slug:       request.Slug,
		HttpMethod: strings.ToLower(request.HttpMethod),
		HttpPath:   request.HttpPath,
		Order:      request.Order,
		ParentId:   request.ParentId,
	}

	u.Create()

	response.Data(c, u)
}

func (uc *AdminPermissionController) Update(c *gin.Context) {
	userModel := adminPermission.Get(c.Param("id"))
	if userModel.ID <= 0 {
		response.Fail(c, "没有找到")
		return
	}

	// 验证
	request := requests.AdminPermissionUpdateRequest{}
	request.ID = userModel.ID
	if ok := requests.ValidateFunc(c, &request, requests.VerityAdminPermissionUpdate); !ok {
		return
	}

	userModel.HttpMethod = strings.ToLower(request.HttpMethod)
	userModel.HttpPath = request.HttpPath
	userModel.Name = request.Name
	userModel.Slug = request.Slug
	userModel.Order = request.Order
	userModel.ParentId = request.ParentId

	userModel.Save()

	response.Data(c, userModel)
}

func (uc *AdminPermissionController) Delete(c *gin.Context) {
	userModel := adminPermission.Get(c.Param("id"))
	if userModel.ID <= 0 {
		response.Fail(c, "没有找到")
		return
	}

	if res := userModel.Delete(); res > 0 {
		response.Success(c)
		return
	}

	response.Fail(c, "删除失败")

}
