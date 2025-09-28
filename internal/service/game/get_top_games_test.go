//go:build unit
// +build unit

package gameservice

import (
	"context"
	"errors"
	"testing"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetTopGames(t *testing.T) {
	t.Run("Internal ошибка", func(t *testing.T) {
		suite := NewSuite()
		filters := dto.GameFilters{ReleaseYear: 2020}
		suite.gameMockRepo.On("GetTopGames", mock.Anything, filters, uint32(10)).Return(([]model.ShortGame)(nil), errors.New("err")).Once()
		games, err := suite.gameService.GetTopGames(context.Background(), filters, uint32(10))
		require.Error(t, err)
		require.Nil(t, games)
	})
	t.Run("Если игр под фильтры не нашлось, ошибки нет", func(t *testing.T) {
		suite := NewSuite()
		filters := dto.GameFilters{ReleaseYear: 2020}
		suite.gameMockRepo.On("GetTopGames", mock.Anything, filters, uint32(10)).Return(([]model.ShortGame)(nil), nil).Once()
		games, err := suite.gameService.GetTopGames(context.Background(), filters, uint32(10))
		require.NoError(t, err)
		require.Nil(t, games)
	})
	t.Run("Успешное получение топа игр", func(t *testing.T) {
		suite := NewSuite()
		filters := dto.GameFilters{ReleaseYear: 2020}
		suite.gameMockRepo.On("GetTopGames", mock.Anything, filters, uint32(10)).Return([]model.ShortGame{{GameID: 1, Title: "qwe", Description: "qe"}}, nil).Once()
		games, err := suite.gameService.GetTopGames(context.Background(), filters, uint32(10))
		require.NoError(t, err)
		require.Equal(t, games, []model.ShortGame{{GameID: 1, Title: "qwe", Description: "qe"}})
	})
}
