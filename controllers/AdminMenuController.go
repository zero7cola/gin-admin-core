package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/model/adminMenu"
	"github.com/zero7cola/gin-admin-core/pkg/response"
	"github.com/zero7cola/gin-admin-core/requests"
)

type AdminMenuController struct {
	BaseAPIController
}

func (uc *AdminMenuController) Index(c *gin.Context) {

	data, pager := adminMenu.Paginate(c, 5)

	response.Data(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

func (uc *AdminMenuController) All(c *gin.Context) {

	menus := adminMenu.All()

	response.Data(c, menus)
}

func (uc *AdminMenuController) Get(c *gin.Context) {

	user := adminMenu.Get(c.Param("id"))

	response.Data(c, user)
}

func (uc *AdminMenuController) Store(c *gin.Context) {
	// 验证
	request := requests.AdminMenuStoreRequest{}
	if ok := requests.ValidateFunc(c, &request, requests.VerityAdminMenuStore); !ok {
		return
	}

	u := adminMenu.AdminMenu{
		Name:     request.Name,
		Icon:     request.Icon,
		Uri:      request.Uri,
		Order:    request.Order,
		ParentId: request.ParentId,
	}

	u.Create()

	response.Data(c, u)
}

func (uc *AdminMenuController) Update(c *gin.Context) {
	userModel := adminMenu.Get(c.Param("id"))
	if userModel.ID <= 0 {
		response.Fail(c, "没有找到")
		return
	}

	// 验证
	request := requests.AdminMenuUpdateRequest{}
	request.ID = userModel.ID
	if ok := requests.ValidateFunc(c, &request, requests.VerityAdminMenuUpdate); !ok {
		return
	}

	userModel.Uri = request.Uri
	userModel.Icon = request.Icon
	userModel.Name = request.Name
	userModel.Order = request.Order
	userModel.ParentId = request.ParentId

	userModel.Save()

	response.Data(c, userModel)
}

func (uc *AdminMenuController) Delete(c *gin.Context) {
	userModel := adminMenu.Get(c.Param("id"))
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
