//go:build integrations

package gamerepo

import "github.com/sariya23/game_service/tests/postgresql"

var (
	db     *postgresql.TestDB
	tables = []string{"game", "tag", "genre", "game_genre", "game_tag"}
)

func init() {
	db = postgresql.NewTestDB()
}
