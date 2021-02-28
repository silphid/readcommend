package db

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DB represents the database and all its tables
type DB struct {
	pool *pgxpool.Pool
}

// Queryer is the abstraction of something able to perform a database query
type Queryer interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

// New creates a new DB object with a connection pool (the
// initial call to create that pool can be cancelled via
// given context). Make sure to call Close() when you're
// done to release connection pool.
func New(ctx context.Context, dbURL string) (*DB, error) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &DB{
		pool: pool,
	}, nil
}

// Query executes given query against database and returns resulting rows.
func (db DB) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return db.pool.Query(ctx, query, args...)
}

// Close releases all resources related to database,
// notably the connection pool.
func (db DB) Close() {
	db.pool.Close()
}

// Check allows to make a very summary sanity check
// with the database, to determine whether it's
// reachable and responsive.
func (db DB) Check(ctx context.Context) error {
	_, err := db.pool.Exec(ctx, "select 1")
	return err
}
