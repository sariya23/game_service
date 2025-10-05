package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	cluster *pgxpool.Pool
}

func MustNewConnection(ctx context.Context, log *slog.Logger, dbURL string) *Database {
	const opearationPlace = "db.MustNewConnection"
	localLog := log.With("operationPlace", opearationPlace)
	ctx, cancel := context.WithTimeout(ctx, time.Second*4)
	defer cancel()
	conn, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		localLog.Error(fmt.Sprintf("%s: cannot connect to db with URL: %s, with error: %v", opearationPlace, dbURL, err))
		panic(fmt.Sprintf("%s: cannot connect to db with URL: %s, with error: %v", opearationPlace, dbURL, err))
	}
	err = conn.Ping(ctx)
	if err != nil {
		localLog.Error(fmt.Sprintf("%s: db is unreachable: %v", opearationPlace, err))
		panic(fmt.Sprintf("%s: db is unreachable: %v", opearationPlace, err))
	}
	localLog.Info("Postgres ready to get connections")
	return &Database{cluster: conn}
}

func (d *Database) Close() {
	d.cluster.Close()
}

func (d *Database) GetPool() *pgxpool.Pool {
	return d.cluster
}
