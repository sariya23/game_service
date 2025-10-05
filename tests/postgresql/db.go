//go:build integrations

package postgresql

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/storage/db"
)

type TestDB struct {
	conn *pgxpool.Pool
}

func NewTestDB() *TestDB {
	cfg := config.MustLoadByPath(filepath.Join("..", "..", "..", "..", "config", "test.env"))
	conn, err := db.NewConnection(
		context.Background(),
		db.GenerateDBUrl(
			cfg.Postgres.PostgresUsername,
			cfg.Postgres.PostgresPassword,
			cfg.Postgres.PostgresHost,
			cfg.Postgres.PostgresPort,
			cfg.Postgres.PostgresDBName,
			cfg.Postgres.SSLMode,
		),
	)
	if err != nil {
		panic(err)
	}
	return &TestDB{conn}
}

func (d *TestDB) Close() {
	d.conn.Close()
}

func (d *TestDB) GetPool() *pgxpool.Pool {
	return d.conn
}

func (d *TestDB) SetUp(ctx context.Context, t *testing.T, tablenames ...string) {
	t.Helper()
	d.Truncate(ctx, tablenames...)
}

func (d *TestDB) TearDown(t *testing.T) {
	t.Helper()
}

func (d *TestDB) Truncate(ctx context.Context, tables ...string) {
	q := fmt.Sprintf("truncate %s", strings.Join(tables, ","))
	if _, err := d.conn.Exec(ctx, q); err != nil {
		panic(err)
	}
}
