package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/outerror"
	gamev2 "github.com/sariya23/proto_api_games/v5/gen/gamev2"
)

func (postgresql PostgreSQL) GameList(ctx context.Context, filters dto.GameFilters, limit uint32) ([]model.ShortGame, error) {
	const operationPlace = "postgresql.GetTopGames"
	log := postgresql.log.With("operationPlave", operationPlace)
	gameGenresQuery := sq.Select(gameGameIDFieldName).From("game")
	gameTagsQuery := sq.Select(gameGameIDFieldName).From("game")

	if t := filters.Tags; len(t) > 0 {
		gameTagsQuery = sq.Select(
			gameTagGameIDFieldName).
			From("tag").
			Join(fmt.Sprintf("game_tag using(%s)", tagTagIDFieldName)).
			Where(sq.Eq{tagTagNameFieldName: t})
	}
	if g := filters.Genres; len(g) > 0 {
		gameGenresQuery = sq.
			Select(gameGenreGameIDFieldName).
			From("genre").
			Join(fmt.Sprintf("game_genre using(%s)", genreGenreIDFieldName)).
			Where(sq.Eq{genreGenreNameFieldName: g})
	}
	tagSQL, tagArgs, _ := gameTagsQuery.ToSql()
	genreSQL, genreArgs, _ := gameGenresQuery.ToSql()

	intersectGameID := fmt.Sprintf("(%s intersect %s)", tagSQL, genreSQL)
	args := append(tagArgs, genreArgs...)

	filteredGameID := sq.Select(
		gameGameIDFieldName,
		gameTitleFieldName,
		gameDescriptionFieldName,
		gameReleaseDateFieldName,
		gameImageURLFieldName,
	).
		From("game").
		Where(sq.Expr(fmt.Sprintf("%s in %s", gameGameIDFieldName, intersectGameID), args...)).
		Where(sq.Eq{gameGameStatusIDFieldName: gamev2.GameStatusType_PUBLISH})

	yearArgs := make([]interface{}, 0, 1)
	if filters.ReleaseYear > 0 {
		filteredGameID = filteredGameID.
			Where(fmt.Sprintf("extract(year from %s)=?", gameReleaseDateFieldName))
		yearArgs = append(yearArgs, filters.ReleaseYear)
	}
	filteredGameID = filteredGameID.
		OrderBy(gameTitleFieldName, gameReleaseDateFieldName).
		Limit(uint64(limit))
	finalSQL, finalArgs, err := filteredGameID.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Error("cannot translate final query to sql string", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	finalArgs = append(finalArgs, yearArgs...)

	var games []model.ShortGame
	gameRows, err := postgresql.connection.Query(ctx, finalSQL, finalArgs...)
	if err != nil {
		log.Error("cannot execute query to get game ids", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer gameRows.Close()
	for gameRows.Next() {
		var game model.ShortGame
		err = gameRows.Scan(
			&game.GameID,
			&game.Title,
			&game.Description,
			&game.ReleaseDate,
			&game.ImageURL,
		)
		if err != nil {
			log.Error("cannot scan game id", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		if gameRows.Err() != nil {
			log.Error("cannot prepare next row", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		games = append(games, game)
	}

	return games, nil
}

func (postgresql PostgreSQL) DaleteGame(ctx context.Context, gameID int64) (*dto.DeletedGame, error) {
	const operationPlace = "postgresql.DeleteGame"
	log := postgresql.log.With("operationPlace", operationPlace)
	deleteGameQuery := fmt.Sprintf(
		"delete from game where %s=$1 returning %s, extract(year from %s), %s",
		gameGameIDFieldName,
		gameGameIDFieldName,
		gameReleaseDateFieldName,
		gameTitleFieldName,
	)
	var deltedGameInfo dto.DeletedGame
	deleteGameRow := postgresql.connection.QueryRow(ctx, deleteGameQuery, gameID)
	err := deleteGameRow.Scan(&deltedGameInfo.GameID, &deltedGameInfo.ReleaseYear, &deltedGameInfo.Title)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("cannot delete game because it is not found", slog.Int("gameID", int(gameID)))
			return nil, fmt.Errorf("%s: %w", operationPlace, outerror.ErrGameNotFound)
		}
		log.Error("cannot delete game", slog.Any("gameID", gameID), slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	log.Info("game deleted successfully", slog.Int("gameID", int(gameID)))
	return &deltedGameInfo, nil
}

func (postgresql PostgreSQL) UpdateGameStatus(ctx context.Context, gameID int64, newStatus gamev2.GameStatusType) error {
	const operationPlace = "postgresql.UpdateGameStatus"
	log := postgresql.log.With("operationPlace", operationPlace)
	queryUpdateStatusQuery := fmt.Sprintf("update game set %s=$1 where %s=$2", gameGameStatusIDFieldName, gameGameIDFieldName)
	_, err := postgresql.connection.Exec(ctx, queryUpdateStatusQuery, newStatus, gameID)
	if err != nil {
		log.Error("cannot update game status", slog.Int64("gameID", gameID), slog.Any("newStatus", newStatus))
		return fmt.Errorf("%s: %w", operationPlace, err)
	}
	return nil
}
