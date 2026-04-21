package controllers

import (
	"time"

	fileModel "github.com/zero7cola/gin-admin-core/model/file"
	"github.com/zero7cola/gin-admin-core/pkg/auth"
	"github.com/zero7cola/gin-admin-core/pkg/response"
	"github.com/zero7cola/gin-admin-core/requests"
	"github.com/zero7cola/gin-admin-core/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type AdminFileController struct {
	BaseAPIController
}

func (uc *AdminFileController) Index(c *gin.Context) {

	data, pager := fileModel.Paginate(c, 5)

	response.Data(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

func (uc *AdminFileController) All(c *gin.Context) {

	menus := fileModel.All()

	response.Data(c, menus)
}

func (uc *AdminFileController) Get(c *gin.Context) {

	user := fileModel.Get(c.Param("id"))

	response.Data(c, user)
}

func (uc *AdminFileController) Store(c *gin.Context) {
	// 验证
	request := requests.AdminFileStoreRequest{}
	if ok := requests.ValidateFunc(c, &request, requests.VerityAdminFileStore); !ok {
		return
	}

	u := fileModel.File{
		Name:         request.Name,
		OriginName:   request.OriginName,
		Size:         request.Size,
		Ext:          request.Ext,
		Type:         request.Type,
		Storage:      request.Storage,
		Url:          request.Url,
		Path:         request.Path,
		LastModified: time.Now(),
		UserId:       cast.ToUint64(auth.CurrentAdminUID(c)),
	}

	u.Create()

	response.Data(c, u)
}

func (uc *AdminFileController) Upload(c *gin.Context) {

	uploadStorage := c.PostForm("uploadStorage")
	obj, err := service.NewFileService(uploadStorage).UploadFile(c)

	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	if obj != nil {
		obj.Create()
	}

	response.Data(c, obj)
}

func (uc *AdminFileController) Update(c *gin.Context) {
	model := fileModel.Get(c.Param("id"))
	if model.ID <= 0 {
		response.Fail(c, "没有找到")
		return
	}

	// 验证
	request := requests.AdminFileUpdateRequest{}
	request.ID = model.ID
	if ok := requests.ValidateFunc(c, &request, requests.VerityAdminFileUpdate); !ok {
		return
	}

	model.Bucket = request.Bucket
	model.Name = request.Name
	model.OriginName = request.OriginName
	model.Path = request.Path
	model.Key = request.Key
	model.Size = request.Size
	model.Ext = request.Ext
	model.Storage = request.Storage
	model.ETag = request.ETag
	model.ContentType = request.ContentType
	model.LastModified = request.LastModified
	model.Url = request.Url
	model.UserId = request.UserId
	model.GroupId = request.GroupId
	model.Type = request.Type

	model.Save()

	response.Data(c, model)
}

func (uc *AdminFileController) Delete(c *gin.Context) {

	model := fileModel.Get(c.Param("id"))
	if model.ID <= 0 {
		response.Fail(c, "没有找到")
		return
	}
	err := service.NewFileService(model.Storage).DeleteFile(cast.ToString(model.ID))

	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c)
}
