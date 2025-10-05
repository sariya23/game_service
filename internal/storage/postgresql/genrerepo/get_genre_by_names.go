package genrerepo

import (
	"context"
	"fmt"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
)

func (gr *GenreRepository) GetGenreByNames(ctx context.Context, genres []string) ([]model.Genre, error) {
	const operationPlace = "postgresql.GetGenres"
	log := gr.log.With("operationPlace", operationPlace)
	getGenresQuery := fmt.Sprintf("select %s, %s from genre where %s=any($1)", GenreGenreIDFieldName, GenreGenreNameFieldName, GenreGenreNameFieldName)
	genreRows, err := gr.conn.GetPool().Query(ctx, getGenresQuery, genres)
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
