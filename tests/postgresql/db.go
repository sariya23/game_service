//go:build integrations

package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"path/filepath"

	"github.com/jackc/pgx/v5"
	apigame "github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/model"
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
	query := "insert into game_genre values ($1, $2)"
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
	query := "insert into game_tag values ($1, $2)"
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

func (d *TestDB) GetTagsByIDs(ctx context.Context, tagIDs []int64) []model.Tag {
	query := "select tag_id, tag_name from tag where tag_id = any($1)"
	rows, err := d.DB.GetPool().Query(ctx, query, tagIDs)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	tags := make([]model.Tag, 0, len(tagIDs))
	for rows.Next() {
		var tag model.Tag
		err = rows.Scan(&tag.TagID, &tag.TagName)
		if err != nil {
			panic(err)
		}
		if rows.Err() != nil {
			panic(err)
		}
		tags = append(tags, tag)
	}
	return tags
}

func (d *TestDB) GetTagsByNames(ctx context.Context, tagNames []string) []model.Tag {
	query := "select tag_id, tag_name from tag where tag_name = any($1)"
	rows, err := d.DB.GetPool().Query(ctx, query, tagNames)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	tags := make([]model.Tag, 0, len(tagNames))
	for rows.Next() {
		var tag model.Tag
		err = rows.Scan(&tag.TagID, &tag.TagName)
		if err != nil {
			panic(err)
		}
		if rows.Err() != nil {
			panic(err)
		}
		tags = append(tags, tag)
	}
	return tags
}

func (d *TestDB) GetTags(ctx context.Context) []model.Tag {
	rows, err := d.DB.GetPool().Query(ctx, "select tag_id, tag_name from tag")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var tags []model.Tag
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.TagID, &tag.TagName); err != nil {
			panic(err)
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return tags
}

func (d *TestDB) GetGenresByIDs(ctx context.Context, genreIDs []int64) []model.Genre {
	query := "select genre_id, genre_name from genre where genre_id = any($1)"
	rows, err := d.DB.GetPool().Query(ctx, query, genreIDs)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	genres := make([]model.Genre, 0, len(genreIDs))
	for rows.Next() {
		var genre model.Genre
		err = rows.Scan(&genre.GenreID, &genre.GenreName)
		if err != nil {
			panic(err)
		}
		if rows.Err() != nil {
			panic(err)
		}
		genres = append(genres, genre)
	}
	return genres
}

func (d *TestDB) GetGenresByNames(ctx context.Context, genreIDs []string) []model.Genre {
	query := "select genre_id, genre_name from genre where genre_name = any($1)"
	rows, err := d.DB.GetPool().Query(ctx, query, genreIDs)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	genres := make([]model.Genre, 0, len(genreIDs))
	for rows.Next() {
		var genre model.Genre
		err = rows.Scan(&genre.GenreID, &genre.GenreName)
		if err != nil {
			panic(err)
		}
		if rows.Err() != nil {
			panic(err)
		}
		genres = append(genres, genre)
	}
	return genres
}

func (d *TestDB) GetGenres(ctx context.Context) []model.Genre {
	rows, err := d.DB.GetPool().Query(ctx, "select genre_id, genre_name from genre")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var genres []model.Genre
	for rows.Next() {
		var genre model.Genre
		if err := rows.Scan(&genre.GenreID, &genre.GenreName); err != nil {
			panic(err)
		}
		genres = append(genres, genre)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return genres
}

func (d *TestDB) GetGameGenreByGameID(ctx context.Context, gameID int64) []model.GameGenre {
	query := "select game_id, genre_id from game_genre where game_id = $1"
	rows, err := d.DB.GetPool().Query(ctx, query, gameID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var gameGenres []model.GameGenre
	for rows.Next() {
		var gameGenre model.GameGenre
		err = rows.Scan(&gameGenre.GameID, &gameGenre.GenreID)
		if err != nil {
			panic(err)
		}
		if rows.Err() != nil {
			panic(err)
		}
		gameGenres = append(gameGenres, gameGenre)
	}
	return gameGenres
}

func (d *TestDB) GetGameTagByGameID(ctx context.Context, gameID int64) []model.GameTag {
	query := "select game_id, tag_id from game_tag where game_id = $1"
	rows, err := d.DB.GetPool().Query(ctx, query, gameID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var gameTags []model.GameTag
	for rows.Next() {
		var gameTag model.GameTag
		err = rows.Scan(&gameTag.GameID, &gameTag.TagID)
		if err != nil {
			panic(err)
		}
		if rows.Err() != nil {
			panic(err)
		}
		gameTags = append(gameTags, gameTag)
	}
	return gameTags
}

func (d *TestDB) GetGameById(ctx context.Context, gameID int64) *model.Game {
	queryGame := "select game_id, title, description, release_date, image_url, game_status_id from game where game_id=$1"
	queryGenre := "select genre_id, genre_name from game join game_genre using(game_id) join genre using(genre_id) where game_id=$1"
	queryTag := "select tag_id, tag_name from game_tag join game using(game_id) join tag using(tag_id) where game_id=$1"
	var game model.Game
	err := d.DB.GetPool().QueryRow(ctx, queryGame, gameID).Scan(
		&game.GameID,
		&game.Title,
		&game.Description,
		&game.ReleaseDate,
		&game.ImageURL,
		&game.GameStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		panic(err)
	}

	rows, err := d.DB.GetPool().Query(ctx, queryGenre, gameID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var genre model.Genre
		err = rows.Scan(&genre.GenreID, &genre.GenreName)
		if err != nil {
			panic(err)
		}
		if rows.Err() != nil {
			panic(err)
		}
		game.Genres = append(game.Genres, genre)
	}

	rows, err = d.DB.GetPool().Query(ctx, queryTag, gameID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tag model.Tag
		err = rows.Scan(&tag.TagID, &tag.TagName)
		if err != nil {
			panic(err)
		}
		if rows.Err() != nil {
			panic(err)
		}
		game.Tags = append(game.Tags, tag)
	}

	return &game
}

func (d *TestDB) UpdateGameStatus(ctx context.Context, gameID int64, newStatus apigame.GameStatusType) {
	query := "update game set game_status_id = $1 where game_id = $2"
	_, err := d.DB.GetPool().Exec(ctx, query, newStatus, gameID)
	if err != nil {
		panic(err)
	}
}
