package minioclient

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	log        *slog.Logger
	client     *minio.Client
	BucketName string
}

func NewMinioClient(
	log *slog.Logger,
	host string,
	port string,
	bucketName string,
	accessKey string,
	secretKey string,
	useSSL bool,
) *Minio {
	client, err := minio.New(fmt.Sprintf("%s:%s", host, port), &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}
	return &Minio{
		log:        log,
		client:     client,
		BucketName: bucketName,
	}
}

func (m Minio) SaveObject(ctx context.Context, name string, data io.Reader) (string, error) {
	panic("impl me")
}

func (m Minio) GetObject(ctx context.Context, name string) (io.Reader, error) {
	panic("impl me")
}

func (m Minio) DeleteObject(ctx context.Context, name string) error {
	panic("imple me")
}

func GameKey(gameTitle string, gameReleaseYear int) string {
	return fmt.Sprintf("%s_%d", gameTitle, gameReleaseYear)
}
