package gameservice

import "github.com/sariya23/game_service/internal/lib/mockslog"

type Suite struct {
	gameService   *GameService
	gameMockRepo  *mockGameReposiroy
	tagMockRepo   *mockTagRepository
	genreMockRepo *mockGenreRepository
	s3Mock        *mockS3Storager
}

func NewSuite() *Suite {
	gameMockRepo := new(mockGameReposiroy)
	tagMockRepo := new(mockTagRepository)
	genreMockRepo := new(mockGenreRepository)
	s3Mock := new(mockS3Storager)
	gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
	return &Suite{gameService: gameService,
		gameMockRepo:  gameMockRepo,
		tagMockRepo:   tagMockRepo,
		genreMockRepo: genreMockRepo,
		s3Mock:        s3Mock,
	}
}
