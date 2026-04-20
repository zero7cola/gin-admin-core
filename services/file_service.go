package service

import (
	"context"
	"errors"
	"github.com/zero7cola/gin-admin-core/config"
	fileModel "github.com/zero7cola/gin-admin-core/model/file"
	"github.com/zero7cola/gin-admin-core/pkg/file"
	"github.com/zero7cola/gin-admin-core/pkg/helpers"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type FileService struct {
	storage file.IStorage
}

func NewFileService(drive ...string) *FileService {

	fileDrive := config.GetString("storage.driver")

	if len(drive) > 0 {
		fileDrive = drive[0]
	}

	fileConfig := file.Config{
		Driver: fileDrive,
		LocalConfig: file.LocalConfig{
			BasePath:      config.GetString("storage.local.path"),
			PublicBaseURL: config.GetString("storage.local.domain"),
		},
		OssConfig: file.OssConfig{
			Region:     config.GetString("storage.oss.region"),
			BucketName: config.GetString("storage.oss.bucket"),
			Key:        config.GetString("storage.oss.key_id"),
			Secret:     config.GetString("storage.oss.key_secret"),
		},
	}
	fileStorage := file.NewStorage(fileConfig)
	return &FileService{
		storage: fileStorage,
	}
}

func (s *FileService) UploadFile(c *gin.Context) (fileModel.File, error) {
	// 从 form-data 获取文件
	fileObj, header, err := c.Request.FormFile("file")
	if err != nil {
		return fileModel.File{}, err
	}
	defer fileObj.Close()

	// 验证大小
	if header.Size > cast.ToInt64(config.Get("storage.size_limit")) {
		return fileModel.File{}, errors.New("超过最大文件大小")
	}

	// 验证后缀
	extLimit := cast.ToStringSlice(config.Get("storage.ext"))
	if helpers.FindElement(extLimit, strings.ToLower(helpers.GetFileExt(header.Filename))) < 0 {
		return fileModel.File{}, errors.New("文件格式不允许 只允许[ " + strings.Join(extLimit, " ") + " ]")
	}

	input := file.PutObjectInput{
		Key:         header.Filename,
		Size:        header.Size,
		ContentType: header.Header.Get("Content-Type"),
		Reader:      fileObj,
		File:        header,
		Meta:        map[string]string{},
	}

	// 上传
	obj, putErr := s.storage.Put(c, input)
	if putErr != nil {
		return fileModel.File{}, putErr
	}

	// 存入数据库
	fileStore := fileModel.File{}
	fileStore.Bucket = obj.Bucket
	fileStore.Name = obj.Name
	fileStore.OriginName = obj.OriginName
	fileStore.Path = obj.Path
	fileStore.Key = obj.Key
	fileStore.Size = obj.Size
	fileStore.Ext = obj.Ext
	fileStore.Storage = obj.Storage
	fileStore.ETag = obj.ETag
	fileStore.ContentType = obj.ContentType
	fileStore.LastModified = obj.LastModified
	fileStore.Url = obj.URL
	fileStore.UserId = 99
	fileStore.GroupId = cast.ToInt(c.DefaultPostForm("group_id", "99"))
	fileStore.Type = cast.ToInt(c.DefaultPostForm("type", "1"))
	fileStore.Create()

	// 组装url
	fileStore.Url = fileStore.GetFileFullUrl()
	return fileStore, nil
}

func (s *FileService) DeleteFile(id string) error {
	fileObj := fileModel.Get(id)
	if fileObj.ID <= 0 {
		return errors.New("没有找到该条记录")
	}
	// 上传
	err := s.storage.Delete(context.Background(), fileObj.Path)
	if err != nil {
		return err
	}

	fileObj.Delete()

	return nil
}
