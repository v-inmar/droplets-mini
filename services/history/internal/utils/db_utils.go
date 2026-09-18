package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"

	// for migration
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Database interface {
	Connect(dsn string) (*sqlx.DB, error)
}

type PostgresDB struct {
	DSN string
}

func NewPostgresDB(dsn string) *PostgresDB {
	return &PostgresDB{
		DSN: dsn,
	}
}

// Connect to postgres with retry
func (db *PostgresDB) Connect(ctx context.Context) (*sqlx.DB, error) {
	timeout := 30 * time.Second
	retryDelay := 5 * time.Second
	connCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		conn, err := sqlx.ConnectContext(connCtx, "postgres", db.DSN)
		if err == nil {
			conn.SetMaxOpenConns(5)
			conn.SetMaxIdleConns(5)
			conn.SetConnMaxLifetime(5 * time.Minute)
			return conn, nil
		}

		select {
		case <-connCtx.Done():
			return nil, fmt.Errorf("connecting to postgres error: %w", connCtx.Err())
		case <-time.After(retryDelay):
		}

	}

}

type Migrate interface {
	Run() error
}

type PostgresMigrate struct {
	src string
	dsn string
}

func NewPostgresMigrate(src string, dsn string) *PostgresMigrate {
	return &PostgresMigrate{
		src: src,
		dsn: dsn,
	}
}

func (pm *PostgresMigrate) Run() error {
	mi, err := migrate.New(
		pm.src,
		pm.dsn,
	)
	if err != nil {
		return err
	}
	defer mi.Close()

	if err := mi.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
