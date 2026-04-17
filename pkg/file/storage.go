package file

import (
	"context"
	"fmt"
	"github.com/zero7cola/gin-admin-core/config"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type PutObjectInput struct {
	Key         string
	Size        int64
	ContentType string
	Reader      io.Reader
	File        *multipart.FileHeader
	// 供直传/分片扩展使用
	Meta map[string]string
}

type ObjectInfo struct {
	Bucket       string
	Key          string
	OriginName   string
	Name         string
	Size         int64
	Storage      string
	ContentType  string
	Ext          string
	Path         string
	ETag         string
	URL          string // 公开或带签名的可访问 URL（可选）
	LastModified time.Time
}

type PresignParams struct {
	Key         string
	Method      string // PUT/GET
	Expiry      time.Duration
	ContentType string
}

type IStorage interface {
	Put(ctx context.Context, in PutObjectInput) (ObjectInfo, error)
	//Get(ctx context.Context, key string) (io.ReadCloser, ObjectInfo, error)
	Delete(ctx context.Context, full string) error
	BackendName() string
}

type Config struct {
	Driver string
	LocalConfig
	OssConfig
}

func NewStorage(cfg Config) IStorage {
	switch cfg.Driver {
	case "local":
		return NewLocalStorage(cfg.LocalConfig)
	case "oss":
		return NewOssStorage(cfg.OssConfig)
	case "other":
		return NewOtherStorage()
	default:
		panic("unsupported storage type")
	}
}

func safeJoin(base, p string) (string, error) {
	clean := filepath.Clean("/" + p)
	// 防目录穿越
	if strings.Contains(clean, "..") {
		return "", fmt.Errorf("invalid path")
	}
	return filepath.Join(base, clean), nil
}

func GetFileStoragePath() string {
	formatted := time.Now().Format("20060102")

	return config.GetString("app.name") + "/" + formatted
}

// 获取文件存储名称(包含完整路径)
func GetFileStorageFullPath(fileName string, isOriginName bool) (string, string) {
	originFileName := fileName
	if !isOriginName {
		fileOriExt := filepath.Ext(fileName) // 获取文件扩展名 这里包含了 .
		//randomNumber := app.GetRandomNumber(16)
		randomNumber := uuid.New().String()
		// fileNameNoExt := fileName[:len(fileName)-len(fileOriExt)] // 文件名称 不含 .和后缀
		originFileName = cast.ToString(randomNumber) + fileOriExt
	}

	objectName := GetFileStoragePath() + "/" + originFileName

	return objectName, originFileName
}
