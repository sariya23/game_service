//go:build integrations

package postgresql

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"path/filepath"

	"github.com/jackc/pgx/v5"
	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/storage/db"
	"github.com/sariya23/game_service/internal/storage/postgresql/gamerepo"
)

type TestDB struct {
	DB *db.Database
}

func NewTestDB() *TestDB {
	cfg := config.MustLoadByPath(filepath.Join("..", "..", "..", "..", "config", "test.env"))
	DB, err := db.NewConnection(
		context.Background(),
		db.GenerateDBUrl(
			cfg.Postgres.PostgresUsername,
			cfg.Postgres.PostgresPassword,
			cfg.Postgres.PostgresHostOuter,
			cfg.Postgres.PostgresPort,
			cfg.Postgres.PostgresDBName,
			cfg.Postgres.SSLMode,
		),
	)
	if err != nil {
		panic(err)
	}
	return &TestDB{DB: DB}
}

func (d *TestDB) SetUp(ctx context.Context, t *testing.T, tablenames ...string) {
	t.Helper()
	d.Truncate(ctx, tablenames...)
}

func (d *TestDB) TearDown(t *testing.T) {
	t.Helper()
}

func (d *TestDB) Truncate(ctx context.Context, tables ...string) {
	q := fmt.Sprintf("truncate %s", strings.Join(tables, ","))
	if _, err := d.DB.GetPool().Exec(ctx, q); err != nil {
		panic(err)
	}
}

func (d *TestDB) InsertGame(ctx context.Context, game dto.AddGameService) int64 {
	query := fmt.Sprintf(
		"insert into game (%s, %s, %s, %s, %s) values ($1, $2, $3, $4, $5) returning game_id;",
		gamerepo.GameTitleFieldName,
		gamerepo.GameDescriptionFieldName,
		gamerepo.GameReleaseDateFieldName,
		gamerepo.GameImageURLFieldName,
		gamerepo.GameGameStatusIDFieldName,
	)
	var id int64
	err := d.DB.GetPool().QueryRow(ctx, query, game.Title, game.Description, game.ReleaseDate, game.ImageURL, 2).Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}

func (d *TestDB) InsertGameGenre(ctx context.Context, gameID int64, genreIDs []int64) {
	query := fmt.Sprintf("insert into game_genre values ($1, $2)")
	batch := &pgx.Batch{}
	for _, genreID := range genreIDs {
		batch.Queue(query, gameID, genreID)
	}
	br := d.DB.GetPool().SendBatch(ctx, batch)
	_, err := br.Exec()
	if err != nil {
		panic(err)
	}
}

func (d *TestDB) InsertGameTag(ctx context.Context, gameID int64, tagIDs []int64) {
	query := fmt.Sprintf("insert into game_tag values ($1, $2)")
	batch := &pgx.Batch{}
	for _, genreID := range tagIDs {
		batch.Queue(query, gameID, genreID)
	}
	br := d.DB.GetPool().SendBatch(ctx, batch)
	_, err := br.Exec()
	if err != nil {
		panic(err)
	}
}
