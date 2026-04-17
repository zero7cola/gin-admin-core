package file

import (
	"context"
	"github.com/zero7cola/gin-admin-core/pkg/helpers"
	"os"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/spf13/cast"
)

type OssConfig struct {
	Region     string
	BucketName string
	Key        string
	Secret     string
}

type OssStorage struct {
	cfg       OssConfig
	ossClient *oss.Client
}

func NewOssStorage(cfg OssConfig) *OssStorage {
	var (
		//region     = "oss-cn-shanghai.aliyuncs.com"
		region = cast.ToString(cfg.Region)
	)
	_ = os.Setenv("OSS_ACCESS_KEY_ID", cast.ToString(cfg.Key))
	_ = os.Setenv("OSS_ACCESS_KEY_SECRET", cast.ToString(cfg.Secret))
	ossCfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)
	client := oss.NewClient(ossCfg)
	return &OssStorage{
		cfg:       cfg,
		ossClient: client,
	}
}

func (l *OssStorage) BackendName() string { return "oss" }

func (l *OssStorage) Put(ctx context.Context, in PutObjectInput) (ObjectInfo, error) {

	// 从请求中获取文件
	fileName := in.File.Filename
	fileSize := in.File.Size
	objectName, nowFileName := GetFileStorageFullPath(fileName, false)

	_, err := l.ossClient.PutObject(ctx, &oss.PutObjectRequest{
		Bucket: oss.Ptr(l.cfg.BucketName),
		Key:    oss.Ptr(objectName),
		Body:   in.Reader,
	})

	if err != nil {
		return ObjectInfo{}, err
	}

	info := ObjectInfo{
		Bucket:       l.cfg.BucketName,
		Key:          nowFileName,
		Name:         nowFileName,
		OriginName:   fileName,
		Path:         objectName,
		Ext:          strings.ToLower(helpers.GetFileExt(in.Key)),
		Storage:      l.BackendName(),
		Size:         fileSize,
		ContentType:  in.ContentType,
		ETag:         "", // 可计算 md5
		LastModified: time.Now(),
	}

	return info, nil
}

func (l *OssStorage) Delete(ctx context.Context, full string) error {

	_, err := l.ossClient.DeleteObject(ctx, &oss.DeleteObjectRequest{
		Bucket: oss.Ptr(l.cfg.BucketName),
		Key:    oss.Ptr(full),
	})

	if err != nil {
		return err
	}

	return nil
}
