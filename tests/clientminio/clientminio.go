package clientminio

import (
	"fmt"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sariya23/game_service/internal/config"
)

type MinioTestClient struct {
	cl         *minio.Client
	BucketName string
}

func NewMinioTestClient() *MinioTestClient {
	cfg := config.MustLoadByPath(filepath.Join("..", "..", "..", "..", "config", "test.env"))
	cl, err := minio.New(
		fmt.Sprintf("%s:%d", cfg.Minio.MinioHostOuter, cfg.Minio.MinioPort),
		&minio.Options{
			Creds:  credentials.NewStaticV4(cfg.Minio.MinioUser, cfg.Minio.MinioPassword, ""),
			Secure: false,
		},
	)
	if err != nil {
		panic(err)
	}
	return &MinioTestClient{cl: cl, BucketName: cfg.Minio.MinioBucket}
}

func (c *MinioTestClient) GetClient() *minio.Client {
	return c.cl
}
