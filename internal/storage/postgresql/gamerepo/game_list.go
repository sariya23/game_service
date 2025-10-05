package gamerepo

import (
	"context"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	gamegenrerepo "github.com/sariya23/game_service/internal/storage/postgresql/game_genre_repo"
	gametagrepo "github.com/sariya23/game_service/internal/storage/postgresql/game_tag_repo"
	"github.com/sariya23/game_service/internal/storage/postgresql/genrerepo"
	"github.com/sariya23/game_service/internal/storage/postgresql/tagrepo"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
)

func (r *GameRepository) GameList(ctx context.Context, filters dto.GameFilters, limit uint32) ([]model.ShortGame, error) {
	const operationPlace = "postgresql.GetTopGames"
	log := r.log.With("operationPlave", operationPlace)
	gameGenresQuery := sq.Select(GameGameIDFieldName).From("game")
	gameTagsQuery := sq.Select(GameGameIDFieldName).From("game")

	if t := filters.Tags; len(t) > 0 {
		gameTagsQuery = sq.Select(
			gametagrepo.GameTagGameIDFieldName).
			From("tag").
			Join(fmt.Sprintf("game_tag using(%s)", tagrepo.TagTagIDFieldName)).
			Where(sq.Eq{tagrepo.TagTagNameFieldName: t})
	}
	if g := filters.Genres; len(g) > 0 {
		gameGenresQuery = sq.
			Select(gamegenrerepo.GameGenreGameIDFieldName).
			From("genre").
			Join(fmt.Sprintf("game_genre using(%s)", genrerepo.GenreGenreIDFieldName)).
			Where(sq.Eq{genrerepo.GenreGenreNameFieldName: g})
	}
	tagSQL, tagArgs, _ := gameTagsQuery.ToSql()
	genreSQL, genreArgs, _ := gameGenresQuery.ToSql()

	intersectGameID := fmt.Sprintf("(%s intersect %s)", tagSQL, genreSQL)
	args := append(tagArgs, genreArgs...)

	filteredGameID := sq.Select(
		GameGameIDFieldName,
		GameTitleFieldName,
		GameDescriptionFieldName,
		GameReleaseDateFieldName,
		GameImageURLFieldName,
	).
		From("game").
		Where(sq.Expr(fmt.Sprintf("%s in %s", GameGameIDFieldName, intersectGameID), args...)).
		Where(sq.Eq{GameGameStatusIDFieldName: gamev2.GameStatusType_PUBLISH})

	yearArgs := make([]interface{}, 0, 1)
	if filters.ReleaseYear > 0 {
		filteredGameID = filteredGameID.
			Where(fmt.Sprintf("extract(year from %s)=?", GameReleaseDateFieldName))
		yearArgs = append(yearArgs, filters.ReleaseYear)
	}
	filteredGameID = filteredGameID.
		OrderBy(GameTitleFieldName, GameReleaseDateFieldName).
		Limit(uint64(limit))
	finalSQL, finalArgs, err := filteredGameID.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Error("cannot translate final query to sql string", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	finalArgs = append(finalArgs, yearArgs...)

	var games []model.ShortGame
	gameRows, err := r.conn.GetPool().Query(ctx, finalSQL, finalArgs...)
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
