package file

import (
	"context"
	"io"
)

type OtherStorage struct {
}

func NewOtherStorage() *OtherStorage { return &OtherStorage{} }

func (l *OtherStorage) BackendName() string { return "local" }

func (l *OtherStorage) Put(ctx context.Context, in PutObjectInput) (ObjectInfo, error) {

	info := ObjectInfo{}
	return info, nil
}

func (l *OtherStorage) Get(ctx context.Context, key string) (io.ReadCloser, ObjectInfo, error) {
	return nil, ObjectInfo{}, nil
}

func (l *OtherStorage) Delete(ctx context.Context, full string) error {
	return nil
}

func (l *OtherStorage) Stat(ctx context.Context, key string) (ObjectInfo, error) {
	return ObjectInfo{}, nil
}

func (l *OtherStorage) Presign(ctx context.Context, p PresignParams) (string, error) {
	return "", nil
}
