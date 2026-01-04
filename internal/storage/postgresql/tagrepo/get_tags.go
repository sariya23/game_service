package tagrepo

import (
	"context"
	"fmt"

	"github.com/sariya23/game_service/internal/lib/logger"
	"github.com/sariya23/game_service/internal/model"
)

func (tr *TagRepository) GetTags(ctx context.Context) ([]model.Tag, error) {
	const operationPlace = "postgresql.GetTags"
	log := tr.log.With("operationPlace", operationPlace)
	log = logger.EnrichRequestID(ctx, log)
	getTagsQuery := fmt.Sprintf("select %s, %s from tag", TagTagIDFieldName, TagTagNameFieldName)
	tagRows, err := tr.conn.GetPool().Query(ctx, getTagsQuery)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get all tags, uncaught error: %v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer tagRows.Close()
	var tagModels []model.Tag
	for tagRows.Next() {
		var modelTag model.Tag
		err = tagRows.Scan(&modelTag.TagID, &modelTag.TagName)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan tags, uncaught error: %v", err))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		tagModels = append(tagModels, modelTag)
	}
	if tagRows.Err() != nil {
		log.Error(fmt.Sprintf("Uncaught error: %v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	return tagModels, nil
}
