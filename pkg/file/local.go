package file

import (
	"context"
	"fmt"
	"github.com/zero7cola/gin-admin-core/pkg/helpers"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LocalConfig struct {
	BasePath      string
	PublicBaseURL string // e.g. http://localhost:8080/static
}

type LocalStorage struct {
	cfg LocalConfig
}

func NewLocalStorage(cfg LocalConfig) *LocalStorage { return &LocalStorage{cfg: cfg} }

func (l *LocalStorage) BackendName() string { return "local" }

func (l *LocalStorage) Put(ctx context.Context, in PutObjectInput) (ObjectInfo, error) {
	full, nowFileName := GetFileStorageFullPath(in.File.Filename, false)
	full = l.cfg.BasePath + "/" + full // 真实的地址

	full = strings.ReplaceAll(full, "\\", "/")
	c := ctx.(*gin.Context)
	err := c.SaveUploadedFile(in.File, full)
	if err != nil {
		return ObjectInfo{}, err
	}

	info := ObjectInfo{
		Bucket:       l.cfg.BasePath,
		Key:          nowFileName,
		Name:         nowFileName,
		OriginName:   in.Key,
		Path:         full,
		Ext:          strings.ToLower(helpers.GetFileExt(in.Key)),
		Storage:      l.BackendName(),
		Size:         in.File.Size,
		ContentType:  in.ContentType,
		ETag:         "", // 可计算 md5
		LastModified: time.Now(),
	}
	//if l.cfg.PublicBaseURL != "" {
	//	u, _ := url.Parse(l.cfg.PublicBaseURL)
	//	u.Path = u.Path + full
	//	info.URL = u.String()
	//}
	return info, nil
}

func (l *LocalStorage) Get(ctx context.Context, key string) (io.ReadCloser, ObjectInfo, error) {
	full, err := safeJoin(l.cfg.BasePath, key)
	if err != nil {
		return nil, ObjectInfo{}, err
	}
	f, err := os.Open(full)
	if err != nil {
		return nil, ObjectInfo{}, err
	}
	st, _ := f.Stat()
	return f, ObjectInfo{Key: key, Size: st.Size(), LastModified: st.ModTime()}, nil
}

func (l *LocalStorage) Delete(ctx context.Context, full string) error {
	return os.Remove(full)
}

func (l *LocalStorage) Stat(ctx context.Context, key string) (ObjectInfo, error) {
	full, err := safeJoin(l.cfg.BasePath, key)
	if err != nil {
		return ObjectInfo{}, err
	}
	st, err := os.Stat(full)
	if err != nil {
		return ObjectInfo{}, err
	}
	return ObjectInfo{Key: key, Size: st.Size(), LastModified: st.ModTime()}, nil
}

func (l *LocalStorage) Presign(ctx context.Context, p PresignParams) (string, error) {
	// 本地存储一般不需要预签名，返回可访问 URL（若配置了）
	if l.cfg.PublicBaseURL == "" {
		return "", fmt.Errorf("no public base url")
	}
	u, _ := url.Parse(l.cfg.PublicBaseURL)
	u.Path = strings.TrimRight(u.Path, "/") + "/" + p.Key
	return u.String(), nil
}
