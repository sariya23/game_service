//go:build integrations

package handler_test

import (
	"context"

	"github.com/sariya23/game_service/tests/postgresql"
)

var (
	dbT      *postgresql.TestDB
	tables   = []string{"game", "game_genre", "game_tag"}
	genreIDs []int64
	tagIDs   []int64
)

func init() {
	ctx := context.Background()
	dbT = postgresql.NewTestDB()
	rows, err := dbT.DB.GetPool().Query(ctx, "select genre_id from genre")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var genre int64
		if err := rows.Scan(&genre); err != nil {
			panic(err)
		}
		genreIDs = append(genreIDs, genre)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	rows, err = dbT.DB.GetPool().Query(ctx, "select tag_id from tag")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tag int64
		if err := rows.Scan(&tag); err != nil {
			panic(err)
		}
		tagIDs = append(tagIDs, tag)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

}
