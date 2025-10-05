package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	Get(ctx context.Context, query string, args ...interface{}) pgx.Row
	Select(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	StartTransaction(ctx context.Context) (pgx.Tx, error)
}

type Database struct {
	cluster *pgxpool.Pool
}

// Get выполняет запрос, который ожидаемо вернет одну строку
func (db Database) Get(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.cluster.QueryRow(ctx, query, args...)
}

// Select выполняет запрос, который ожидаемое вернет несколько строк
func (db Database) Select(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return db.cluster.Query(ctx, query, args...)
}

// Exec выполняте запрос, который ожидаемо не вернет строк
func (db Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.cluster.Exec(ctx, query, args...)
}

func (db Database) StartTransaction(ctx context.Context) (pgx.Tx, error) {
	return db.cluster.Begin(ctx)
}
