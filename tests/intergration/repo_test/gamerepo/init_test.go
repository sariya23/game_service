//go:build integrations

package gamerepo

import (
	"context"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/tests/postgresql"
)

var (
	db     *postgresql.TestDB
	tables = []string{"game", "game_genre", "game_tag"}
	genres []model.Genre
	tags   []model.Tag
)

func init() {
	ctx := context.Background()
	db = postgresql.NewTestDB()
	rows, err := db.GetPool().Query(ctx, "select genre_name from genre")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var genre string
		if err := rows.Scan(&genre); err != nil {
			panic(err)
		}
		genres = append(genres, genre)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	rows, err = db.GetPool().Query(ctx, "select tag_name from tag")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			panic(err)
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

}
