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

func MustPrepareMinio(
	ctx context.Context,
	log *slog.Logger,
	host string,
	port int,
	bucketName string,
	accessKey string,
	secretKey string,
	useSSL bool,
) *Minio {
	min, err := newMinioClient(log, host, port, bucketName, accessKey, secretKey, useSSL)
	if err != nil {
		panic(err)
	}
	err = min.createBucket(ctx)
	if err != nil {
		panic(err)
	}
	return min
}

func newMinioClient(
	log *slog.Logger,
	host string,
	port int,
	bucketName string,
	accessKey string,
	secretKey string,
	useSSL bool,
) (*Minio, error) {
	const operationPlace = "minioclient.NewMinioClient"
	client, err := minio.New(fmt.Sprintf("%s:%d", host, port), &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	return &Minio{
		log:        log,
		client:     client,
		BucketName: bucketName,
	}, nil
}

func (m Minio) createBucket(ctx context.Context) error {
	const operationPlace = "minioclient.CreateBucket"
	log := m.log.With("operationPlace", operationPlace)
	err := m.client.MakeBucket(ctx, m.BucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := m.client.BucketExists(ctx, m.BucketName)
		if errBucketExists == nil && exists {
			log.Info("bucket already exists", slog.String("name", m.BucketName))
			return nil
		}
		return fmt.Errorf("%s: %w", operationPlace, err)
	}
	return nil
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
