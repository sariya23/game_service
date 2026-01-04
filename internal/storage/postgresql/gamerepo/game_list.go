package gamerepo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	gamegenrerepo "github.com/sariya23/game_service/internal/storage/postgresql/game_genre_repo"
	gametagrepo "github.com/sariya23/game_service/internal/storage/postgresql/game_tag_repo"
	"github.com/sariya23/game_service/internal/storage/postgresql/genrerepo"
	"github.com/sariya23/game_service/internal/storage/postgresql/tagrepo"
)

func (gr *GameRepository) GameList(ctx context.Context, filters dto.GameFilters, limit uint32) ([]model.ShortGame, error) {
	const operationPlace = "postgresql.GetTopGames"
	log := gr.log.With("operationPlave", operationPlace)
	baseQuery := fmt.Sprintf("select %s, %s, %s, %s, %s from game where true",
		GameGameIDFieldName,
		GameTitleFieldName,
		GameDescriptionFieldName,
		GameReleaseDateFieldName,
		GameImageURLFieldName,
	)
	args := []interface{}{}
	if len(filters.Tags) > 0 {
		args = append(args, filters.Tags)
		baseQuery = baseQuery + fmt.Sprintf(" and %s in (select %s from game_tag join tag using(%s) where %s=any($%d))",
			GameGameIDFieldName,
			gametagrepo.GameTagGameIDFieldName,
			gametagrepo.GameTagTagIDFieldName,
			tagrepo.TagTagNameFieldName,
			len(args))
	}
	if len(filters.Genres) > 0 {
		args = append(args, filters.Genres)
		baseQuery = baseQuery + fmt.Sprintf(" and %s in (select %s from game_genre join genre using(%s) where %s=any($%d))",
			GameGameIDFieldName,
			gamegenrerepo.GameGenreGameIDFieldName,
			gamegenrerepo.GameGenreGenreIDFieldName,
			genrerepo.GenreGenreNameFieldName,
			len(args))
	}
	if filters.ReleaseYear > 0 {
		args = append(args, filters.ReleaseYear)
		baseQuery = baseQuery + fmt.Sprintf(" and extract(year from %s) = $%d", GameReleaseDateFieldName, len(args))
	}
	baseQuery = baseQuery + fmt.Sprintf(" order by %s, %s limit %d",
		GameTitleFieldName,
		GameReleaseDateFieldName,
		limit,
	)
	var games []model.ShortGame
	log.Info(baseQuery)
	log.Info(fmt.Sprintf("%v", args))
	gameRows, err := gr.conn.GetPool().Query(ctx, baseQuery, args...)
	if err != nil {
		log.Error("cannot execute query to get games", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer gameRows.Close()
	for gameRows.Next() {
		var gameDB dto.ShortGameDB
		err = gameRows.Scan(
			&gameDB.GameID,
			&gameDB.Title,
			&gameDB.Description,
			&gameDB.ReleaseDate,
			&gameDB.ImageKey,
		)
		if err != nil {
			log.Error("cannot scan game id", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		if gameRows.Err() != nil {
			log.Error("cannot prepare next row", slog.String("err", gameRows.Err().Error()))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		game := gameDB.ToDomain()
		games = append(games, game)
	}
	return games, nil
}
