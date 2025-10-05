package genrerepo

import (
	"context"
	"fmt"

	"github.com/sariya23/game_service/internal/model"
)

// GetGenres возвращает все жанры.
func (gr *GenreRepository) GetGenres(ctx context.Context) ([]model.Genre, error) {
	const operationPlace = "postgresql.GetGenres"
	log := gr.log.With("operationPlace", operationPlace)
	getGenreQuery := fmt.Sprintf("select %s, %s from genre", GenreGenreIDFieldName, GenreGenreNameFieldName)
	genreRows, err := gr.conn.GetPool().Query(ctx, getGenreQuery)
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
