package tagrepo

import (
	"context"
	"fmt"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
)

func (tr *TagRepository) GetTagByNames(ctx context.Context, tags []string) ([]model.Tag, error) {
	const operationPlace = "postgresql.GetTags"
	log := tr.log.With("operationPlace", operationPlace)
	getTagsQuery := fmt.Sprintf("select %s, %s from tag where %s=any($1)", TagTagIDFieldName, TagTagNameFieldName, TagTagNameFieldName)
	tagRows, err := tr.conn.GetPool().Query(ctx, getTagsQuery, tags)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get tags from request, uncaught error: %v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer tagRows.Close()
	tagModels := make([]model.Tag, 0, len(tags))
	for tagRows.Next() {
		var modelTag model.Tag
		err = tagRows.Scan(&modelTag.TagID, &modelTag.TagName)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan tags from request, uncaught error: %v", err))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		tagModels = append(tagModels, modelTag)
	}
	if tagRows.Err() != nil {
		log.Error(fmt.Sprintf("Uncaught error: %v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	if len(tagModels) != len(tags) {
		return nil, fmt.Errorf("%s: %w", operationPlace, outerror.ErrTagNotFound)
	}
	return tagModels, nil
}
