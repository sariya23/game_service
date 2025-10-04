package postgresql

import (
	"context"
	"fmt"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
)

// GetTagByNames возвращает тэги по переданны именам.
func (postgresql PostgreSQL) GetTagByNames(ctx context.Context, tags []string) ([]*model.Tag, error) {
	const operationPlace = "postgresql.GetTags"
	log := postgresql.log.With("operationPlace", operationPlace)
	getTagsQuery := fmt.Sprintf("select %s, %s from tag where %s=any($1)", tagTagIDFieldName, tagTagNameFieldName, tagTagNameFieldName)
	tagRows, err := postgresql.connection.Query(ctx, getTagsQuery, tags)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get tags from request, uncaught error: %v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer tagRows.Close()
	tagModels := make([]*model.Tag, 0, len(tags))
	for tagRows.Next() {
		var modelTag model.Tag
		err = tagRows.Scan(&modelTag.TagID, &modelTag.TagName)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan tags from request, uncaught error: %v", err))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		tagModels = append(tagModels, &modelTag)
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

// GetTags возвращает все тэги.
func (postgresql PostgreSQL) GetTags(ctx context.Context) ([]*model.Tag, error) {
	const operationPlace = "postgresql.GetTags"
	log := postgresql.log.With("operationPlace", operationPlace)
	getTagsQuery := fmt.Sprintf("select %s, %s from tag", tagTagIDFieldName, tagTagNameFieldName)
	tagRows, err := postgresql.connection.Query(ctx, getTagsQuery)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get all tags, uncaught error: %v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer tagRows.Close()
	var tagModels []*model.Tag
	for tagRows.Next() {
		var modelTag model.Tag
		err = tagRows.Scan(&modelTag.TagID, &modelTag.TagName)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan tags, uncaught error: %v", err))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		tagModels = append(tagModels, &modelTag)
	}
	if tagRows.Err() != nil {
		log.Error(fmt.Sprintf("Uncaught error: %v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	return tagModels, nil
}
