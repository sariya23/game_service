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

func (postgresql PostgreSQL) GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (*model.Game, error) {
	const operationPlace = "postgresql.GetGameByTitleAndReleaseYear"
	log := postgresql.log.With("operationPlace", operationPlace).With("title", title).With("releaseYear", releaseYear)
	getGameQuery := fmt.Sprintf("select %s, %s, %s, %s, %s, %s from game where %s=$1 and extract(year from %s)=$2",
		gameGameIDFieldName,
		gameTitleFieldName,
		gameDescriptionFieldName,
		gameReleaseDateFieldName,
		gameImageURLFieldName,
		gameGameStatusIDFieldName,
		gameTitleFieldName,
		gameReleaseDateFieldName,
	)
	getGameGenresQuery := fmt.Sprintf(`
	select %s, %s
	from game join game_genre using(%s) join genre using(%s)
	where %s=$1`,
		genreGenreNameFieldName,
		genreGenreIDFieldName,
		gameGenreGameIDFieldName,
		genreGenreIDFieldName,
		gameGameIDFieldName,
	)
	getGameTagsQuery := fmt.Sprintf(`
	select %s, %s
	from game join game_tag using(%s) join tag using(%s)
	where %s=$1`,
		tagTagNameFieldName,
		tagTagIDFieldName,
		gameTagGameIDFieldName,
		tagTagIDFieldName,
		gameGameIDFieldName,
	)
	var game model.Game
	gameRow := postgresql.connection.QueryRow(ctx, getGameQuery, title, releaseYear)
	err := gameRow.Scan(
		&game.GameID,
		&game.Title,
		&game.Description,
		&game.ReleaseDate,
		&game.ImageURL,
		&game.GameStatus,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("No game with this data")
			return nil, outerror.ErrGameNotFound
		}
		log.Error(fmt.Sprintf("cannot get game, unexpected error = %v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	genreRows, err := postgresql.connection.Query(ctx, getGameGenresQuery, game.GameID)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get genres, uncaught error: %v", err), slog.Int64("gameID", game.GameID))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer genreRows.Close()
	for genreRows.Next() {
		var gameGenre model.Genre
		err = genreRows.Scan(&gameGenre.GenreName, &gameGenre.GenreID)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan game genres, uncaught error: %v", err), slog.Int64("gameID", game.GameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		if genreRows.Err() != nil {
			log.Error("cannot prepare next row", slog.String("err", err.Error()), slog.Int64("gameID", game.GameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		game.Genres = append(game.Genres, gameGenre)
	}

	tagRows, err := postgresql.connection.Query(ctx, getGameTagsQuery, game.GameID)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get tags, uncaught error: %v", err), slog.Int64("gameID", game.GameID))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer tagRows.Close()
	for tagRows.Next() {
		var gameTag model.Tag
		err = tagRows.Scan(&gameTag.TagName, &gameTag.TagID)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan game tags, uncaught error: %v", err), slog.Int64("gameID", game.GameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		if tagRows.Err() != nil {
			log.Error("cannot prepare next row", slog.String("err", err.Error()), slog.Int64("gameID", game.GameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		game.Tags = append(game.Tags, gameTag)
	}
	return &game, nil
}

func (postgresql PostgreSQL) GetGameByID(ctx context.Context, gameID int64) (*model.Game, error) {
	const operationPlace = "postgresql.GetGameByID"
	log := postgresql.log.With("operationPlace", operationPlace)
	getGameMainInfoQuery := fmt.Sprintf(
		"select %s, %s, %s, %s, %s, %s from game where %s=$1",
		gameGameIDFieldName,
		gameTitleFieldName,
		gameDescriptionFieldName,
		gameReleaseDateFieldName,
		gameImageURLFieldName,
		gameGameStatusIDFieldName,
		gameGameIDFieldName,
	)
	getGameGenresQuery := fmt.Sprintf(`
	select %s, %s
	from game join game_genre using(%s) join genre using(%s)
	where %s=$1`,
		genreGenreNameFieldName,
		genreGenreIDFieldName,
		gameGenreGameIDFieldName,
		genreGenreIDFieldName,
		gameGameIDFieldName,
	)
	getGameTagsQuery := fmt.Sprintf(`
	select %s, %s 
	from game join game_tag using(%s) join tag using(%s)
	where %s=$1`,
		tagTagNameFieldName,
		tagTagIDFieldName,
		gameTagGameIDFieldName,
		tagTagIDFieldName,
		gameGameIDFieldName,
	)
	var game model.Game
	gameRow := postgresql.connection.QueryRow(ctx, getGameMainInfoQuery, gameID)
	err := gameRow.Scan(
		&game.GameID,
		&game.Title,
		&game.Description,
		&game.ReleaseDate,
		&game.ImageURL,
		&game.GameStatus,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("game does not exists", slog.Int64("gameID", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, outerror.ErrGameNotFound)
		} else {
			log.Error(fmt.Sprintf("Uncaught error: %v", err))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}
	genreRows, err := postgresql.connection.Query(ctx, getGameGenresQuery, gameID)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get genres, uncaught error: %v", err), slog.Int64("gameID", gameID))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer genreRows.Close()
	for genreRows.Next() {
		var gameGenre model.Genre
		err = genreRows.Scan(&gameGenre.GenreName, &gameGenre.GenreID)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan game genres, uncaught error: %v", err), slog.Int64("gameID", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		if genreRows.Err() != nil {
			log.Error("cannot prepare next row", slog.String("err", err.Error()), slog.Int64("gameID", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		game.Genres = append(game.Genres, gameGenre)
	}

	tagRows, err := postgresql.connection.Query(ctx, getGameTagsQuery, gameID)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get tags, uncaught error: %v", err), slog.Int64("gameID", gameID))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer tagRows.Close()
	for tagRows.Next() {
		var gameTag model.Tag
		err = tagRows.Scan(&gameTag.TagName, &gameTag.TagID)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan game tags, uncaught error: %v", err), slog.Int64("gameID", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		if tagRows.Err() != nil {
			log.Error("cannot prepare next row", slog.String("err", err.Error()), slog.Int64("gameID", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		game.Tags = append(game.Tags, gameTag)
	}

	return &game, nil
}

func (postgresql PostgreSQL) SaveGame(ctx context.Context, game model.Game) (int64, error) {
	const operationPlace = "postgresql.SaveGame"
	log := postgresql.log.With("operationPlace", operationPlace)
	saveGameArgs := pgx.NamedArgs{
		"title":        game.Title,
		"description":  game.Description,
		"release_date": game.ReleaseDate,
		"image_url":    game.ImageURL,
	}
	saveMainGameInfoQuery := fmt.Sprintf(`
		insert into game (%s, %s, %s, %s) values 
		(@title, @description, @release_date, @image_url)
		returning %s
	`, gameTitleFieldName, gameDescriptionFieldName, gameReleaseDateFieldName, gameImageURLFieldName, gameGameIDFieldName)
	addTagsForGameQuery := "insert into game_tag values ($1, $2)"
	addGenresForGameQuery := "insert into game_genre values ($1, $2)"
	genreIDs := make([]int, 0, len(game.Genres))
	tagIDs := make([]int, 0, len(game.Tags))

	if len(game.Genres) != 0 {
		for _, g := range game.Genres {
			genreIDs = append(genreIDs, int(g.GenreID))
		}
	}

	if len(game.Tags) != 0 {
		for _, t := range game.Tags {
			tagIDs = append(tagIDs, int(t.TagID))
		}
	}

	tx, err := postgresql.connection.Begin(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("cannot start transaction, unexpected error = %v", err))
		return 0, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer tx.Rollback(ctx)

	var savedGameID int64
	saveGameRow := tx.QueryRow(ctx, saveMainGameInfoQuery, saveGameArgs)
	err = saveGameRow.Scan(&savedGameID)
	if err != nil {
		log.Error(fmt.Sprintf("cannot save game, unexpected error = %v", err))
		return 0, fmt.Errorf("%s: %w", operationPlace, err)
	}

	for _, tagID := range tagIDs {
		_, err = tx.Exec(ctx, addTagsForGameQuery, savedGameID, tagID)
		if err != nil {
			log.Error(fmt.Sprintf("cannot link tag with game, unexpected error = %v", err), slog.Int("tagID", tagID), slog.Int("gameID", int(savedGameID)))
			return 0, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}

	for _, genreID := range genreIDs {
		_, err = tx.Exec(ctx, addGenresForGameQuery, savedGameID, genreID)
		if err != nil {
			log.Error(fmt.Sprintf("cannot link tag with game, unexpected error = %v", err), slog.Int("genreID", genreID), slog.Int("gameID", int(savedGameID)))
			return 0, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("cannot commit, err = %v", err))
		return 0, fmt.Errorf("%s: %w", operationPlace, err)
	}
	return savedGameID, nil
}

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
