//go:build integrations

package game_test

import (
	"github.com/sariya23/game_service/tests/clientminio"
	"github.com/sariya23/game_service/tests/postgresql"
)

var (
	dbT    *postgresql.TestDB
	minioT *clientminio.MinioTestClient
	tables = []string{"game", "game_genre", "game_tag"}
)

func init() {
	dbT = postgresql.NewTestDB()
	minioT = clientminio.NewMinioTestClient()
}
