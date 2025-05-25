package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	pool *pgxpool.Pool
}

func ConnectToPostgres() (*PostgresDB, error) {
	config, err := pgxpool.ParseConfig("postgres://postgres:example@localhost:5432/postgres")

	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	return &PostgresDB{pool: pool}, nil
}
