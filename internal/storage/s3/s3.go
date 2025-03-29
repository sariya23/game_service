package s3

import (
	"context"
	"fmt"
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

func (s3 S3Storage) Save(ctx context.Context, data io.Reader, key string) (string, error) {
	panic("impl me")
}

func (s3 S3Storage) Get(ctx context.Context, bucket, key string) io.Reader {
	panic("impl me")
}

func CreateGameKey(gameTitle string, gameReleaseYear int) string {
	return fmt.Sprintf("%s_%d_image", gameTitle, gameReleaseYear)
}
