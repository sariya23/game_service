package minioclient

import (
	"context"
	"io"
	"log/slog"
)

type Minio struct {
	log        *slog.Logger
	BucketName string
}

func NewMinioClient(log *slog.Logger, bucketName string) *Minio {
	return &Minio{log: log, BucketName: bucketName}
}

func (m *Minio) SaveObject(ctx context.Context, name string, data io.Reader) (string, error) {
	panic("impl me")
}

func (m *Minio) GetObject(ctx context.Context, name string) (io.Reader, error) {
	panic("impl me")
}
