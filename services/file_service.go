package service

import (
	"context"
	"errors"
	"strings"

	"github.com/zero7cola/gin-admin-core/setting"

	fileModel "github.com/zero7cola/gin-admin-core/model/file"
	"github.com/zero7cola/gin-admin-core/pkg/file"
	"github.com/zero7cola/gin-admin-core/pkg/helpers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type FileService struct {
	storage file.IStorage
}

func NewFileService(drive ...string) *FileService {

	fileDrive := setting.GlobalSetting.Storage.Driver

	if len(drive) > 0 {
		fileDrive = drive[0]
	}

	fileConfig := file.Config{
		Driver: fileDrive,
		LocalConfig: file.LocalConfig{
			BasePath:      setting.GlobalSetting.Storage.Local.Path,
			PublicBaseURL: setting.GlobalSetting.Storage.Local.Domain,
		},
		OssConfig: file.OssConfig{
			Region:     setting.GlobalSetting.Storage.Oss.Region,
			BucketName: setting.GlobalSetting.Storage.Oss.Bucket,
			Key:        setting.GlobalSetting.Storage.Oss.KeyId,
			Secret:     setting.GlobalSetting.Storage.Oss.KeySecret,
		},
	}
	fileStorage := file.NewStorage(fileConfig)
	return &FileService{
		storage: fileStorage,
	}
}

func (s *FileService) UploadFile(c *gin.Context) (*fileModel.File, error) {
	// 从 form-data 获取文件
	fileObj, header, err := c.Request.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer fileObj.Close()

	// 验证大小
	if header.Size > cast.ToInt64(setting.GlobalSetting.Storage.SizeLimit) {
		return nil, errors.New("超过最大文件大小")
	}

	// 验证后缀
	extLimit := cast.ToStringSlice(setting.GlobalSetting.Storage.Ext)
	if helpers.FindElement(extLimit, strings.ToLower(helpers.GetFileExt(header.Filename))) < 0 {
		return nil, errors.New("文件格式不允许 只允许[ " + strings.Join(extLimit, " ") + " ]")
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
		return nil, putErr
	}

	// 存入数据库
	fileStore := &fileModel.File{}
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

	// 组装url
	fileStore.FullUrl = fileStore.GetFileFullUrl()
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
