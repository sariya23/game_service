package postgresql

import (
	"context"
	"fmt"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
)

func (postgresql PostgreSQL) GetGenreByNames(ctx context.Context, genres []string) ([]model.Genre, error) {
	const operationPlace = "postgresql.GetGenres"
	log := postgresql.log.With("operationPlace", operationPlace)
	getGenresQuery := fmt.Sprintf("select %s, %s from genre where %s=any($1)", genreGenreIDFieldName, genreGenreNameFieldName, genreGenreNameFieldName)
	genreRows, err := postgresql.connection.Query(ctx, getGenresQuery, genres)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get genres, uncaught error: %v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer genreRows.Close()
	genreModels := make([]model.Genre, 0, len(genres))
	for genreRows.Next() {
		var modelGenre model.Genre
		err = genreRows.Scan(&modelGenre.GenreID, &modelGenre.GenreName)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan game tags, uncaught error: %v", err))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		genreModels = append(genreModels, modelGenre)
	}
	if genreRows.Err() != nil {
		log.Error(fmt.Sprintf("Uncaught error: %v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	if len(genreModels) != len(genres) {
		return nil, fmt.Errorf("%s: %w", operationPlace, outerror.ErrGenreNotFound)
	}
	return genreModels, nil
}

// GetGenres возвращает все жанры.
func (postgresql PostgreSQL) GetGenres(ctx context.Context) ([]model.Genre, error) {
	const operationPlace = "postgresql.GetGenres"
	log := postgresql.log.With("operationPlace", operationPlace)
	getGenreQuery := fmt.Sprintf("select %s, %s from genre", genreGenreIDFieldName, genreGenreNameFieldName)
	genreRows, err := postgresql.connection.Query(ctx, getGenreQuery)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get all tags, uncaught error: %v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer genreRows.Close()
	var genreModels []model.Genre
	for genreRows.Next() {
		var modelGenre model.Genre
		err = genreRows.Scan(&modelGenre.GenreID, &modelGenre.GenreName)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan tags, uncaught error: %v", err))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		genreModels = append(genreModels, modelGenre)
	}
	if genreRows.Err() != nil {
		log.Error(fmt.Sprintf("Uncaught error: %v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	return genreModels, nil
}
