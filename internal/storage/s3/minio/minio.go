package minioclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sariya23/game_service/internal/config"
)

type Minio struct {
	log        *slog.Logger
	client     *minio.Client
	BucketName string
}

func MustPrepareMinio(
	ctx context.Context,
	log *slog.Logger,
	minioConfig *config.Minio,
	useSSL bool,
) *Minio {
	const operationPlace = "minioclient.MustPrepareMinio"
	innerLog := log.With("operationPlace", operationPlace)
	minioClient, err := newMinioClient(log,
		minioConfig.MinioHostInner,
		minioConfig.MinioPort,
		minioConfig.MinioBucket,
		minioConfig.MinioUser,
		minioConfig.MinioPassword,
		useSSL,
	)
	if err != nil {
		panic(err)
	}
	err = minioClient.createBucket(ctx)
	if err != nil {
		panic(err)
	}
	innerLog.Info("Minio ready to get connections")
	return minioClient
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
	logInner := log.With("operationPlace", operationPlace)
	client, err := minio.New(fmt.Sprintf("%s:%d", host, port), &minio.Options{
		Creds:        credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure:       useSSL,
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	logInner.Info("Minio client ready to get connections")
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
	log.Info("Bucket created")
	return nil
}

func (m Minio) SaveObject(ctx context.Context, name string, data io.Reader) (string, error) {
	const operationPlace = "minioclient.SaveObject"
	log := m.log.With("operationPlace", operationPlace)
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(data)
	if err != nil {
		log.Error(fmt.Sprintf("cannot read data from reader; err=%v", err))
		return "", fmt.Errorf("%s: %w", operationPlace, err)
	}
	if buf.Len() == 0 {
		log.Error("empty data provided")
		return "", fmt.Errorf("%s: empty data", operationPlace)
	}
	reader := bytes.NewReader(buf.Bytes())
	info, err := m.client.PutObject(ctx, m.BucketName, name, reader, int64(buf.Len()), minio.PutObjectOptions{})
	if err != nil {
		log.Error(fmt.Sprintf("cannot save object in s3; err=%v", err))
		return "", fmt.Errorf("%s: %w", operationPlace, err)
	}

	return info.Key, nil
}

func (m Minio) DeleteObject(ctx context.Context, name string) error {
	const operationPlace = "minioclient.GetObject"
	log := m.log.With("operationPlace", operationPlace)
	err := m.client.RemoveObject(ctx, m.BucketName, name, minio.RemoveObjectOptions{})
	if err != nil {
		log.Error(fmt.Sprintf("unexpected error; err=%v", err))
		return fmt.Errorf("%s: %w", operationPlace, err)
	}
	return nil
}

func GameKey(gameTitle string, gameReleaseYear int) string {
	return fmt.Sprintf("%s_%d", gameTitle, gameReleaseYear)
}
