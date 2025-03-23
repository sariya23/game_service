package s3

import (
	"context"
	"io"
	"log/slog"
)

type S3Storage struct {
	log *slog.Logger
}

func NewS3Storage(log *slog.Logger) *S3Storage {
	return &S3Storage{
		log: log,
	}
}

func (s3 S3Storage) Save(ctx context.Context, data io.Reader, key string) error {
	panic("impl me")
}

func (s3 S3Storage) Get(ctx context.Context, bucket, key string) io.Reader {
	panic("impl me")
}

func GetImageURL(key string) string {
	return ""
}
