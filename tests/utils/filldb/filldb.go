//go:build integrations

package filldb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/storage/postgresql/gamerepo"
	"github.com/sariya23/game_service/tests/postgresql"
)

func InsertGame(ctx context.Context, db *postgresql.TestDB, game dto.AddGameService) int64 {
	query := fmt.Sprintf(
		"insert into game (%s, %s, %s, %s, %s) values ($1, $2, $3, $4, $5) returning id;",
		gamerepo.GameTitleFieldName,
		gamerepo.GameDescriptionFieldName,
		gamerepo.GameReleaseDateFieldName,
		gamerepo.GameImageURLFieldName,
		gamerepo.GameGameStatusIDFieldName,
	)
	var id int64
	err := db.DB.GetPool().QueryRow(ctx, query, game.Title, game.Description, game.ReleaseDate, game.ImageURL, 2).Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}

func InsertGameGenre(ctx context.Context, db *postgresql.TestDB, gameID int64, genreIDs []int64) {
	query := fmt.Sprintf("insert into game_genre values ($1, $2)")
	batch := &pgx.Batch{}
	for _, genreID := range genreIDs {
		batch.Queue(query, gameID, genreID)
	}
	br := db.DB.GetPool().SendBatch(ctx, batch)
	_, err := br.Exec()
	if err != nil {
		panic(err)
	}
}

func InsertGameTag(ctx context.Context, db *postgresql.TestDB, gameID int64, tagIDs []int64) {
	query := fmt.Sprintf("insert into game_tag values ($1, $2)")
	batch := &pgx.Batch{}
	for _, genreID := range tagIDs {
		batch.Queue(query, gameID, genreID)
	}
	br := db.DB.GetPool().SendBatch(ctx, batch)
	_, err := br.Exec()
	if err != nil {
		panic(err)
	}
}
