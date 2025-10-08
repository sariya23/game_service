//go:build integrations

package game_test

import (
	"context"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/tests/clientminio"
	"github.com/sariya23/game_service/tests/postgresql"
)

var (
	dbT    *postgresql.TestDB
	minioT *clientminio.MinioTestClient
	tables = []string{"game", "game_genre", "game_tag"}
	genres []model.Genre
	tags   []model.Tag
)

func init() {
	ctx := context.Background()
	dbT = postgresql.NewTestDB()
	minioT = clientminio.NewMinioTestClient()
	rows, err := dbT.DB.GetPool().Query(ctx, "select genre_id, genre_name from genre")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var genre model.Genre
		if err := rows.Scan(&genre.GenreID, &genre.GenreName); err != nil {
			panic(err)
		}
		genres = append(genres, genre)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	rows, err = dbT.DB.GetPool().Query(ctx, "select tag_id, tag_name from tag")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.TagID, &tag.TagName); err != nil {
			panic(err)
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

}
