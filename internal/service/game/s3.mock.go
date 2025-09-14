package gameservice

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/mock"
)

type mockS3Storager struct {
	mock.Mock
}

func (m *mockS3Storager) SaveObject(ctx context.Context, name string, data io.Reader) (string, error) {
	args := m.Called(ctx, name, data)
	return args.Get(0).(string), args.Error(1)
}

func (m *mockS3Storager) GetObject(ctx context.Context, name string) (*minio.Object, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*minio.Object), args.Error(1)
}

func (m *mockS3Storager) DeleteObject(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}
