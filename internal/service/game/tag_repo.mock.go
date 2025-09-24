package gameservice

import (
	"context"

	"github.com/sariya23/game_service/internal/model"
	"github.com/stretchr/testify/mock"
)

type mockTagRepository struct {
	mock.Mock
}

func (m *mockTagRepository) GetTagByNames(ctx context.Context, tags []string) ([]model.Tag, error) {
	args := m.Called(ctx, tags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Tag), args.Error(1)
}
