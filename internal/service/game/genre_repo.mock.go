package gameservice

import (
	"context"

	"github.com/sariya23/game_service/internal/model"
	"github.com/stretchr/testify/mock"
)

type mockGenreRepository struct {
	mock.Mock
}

func (m *mockGenreRepository) GetGenreByNames(ctx context.Context, genres []string) ([]model.Genre, error) {
	args := m.Called(ctx, genres)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Genre), args.Error(1)
}
