package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrNotFound     = errors.New("record not found")
	ErrConflict     = errors.New("record already exists")
	ErrInvalidInput = errors.New("invalid input")
)

type Store struct {
	DB *sql.DB
}

func Open(ctx context.Context, dsn string) (*Store, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	return &Store{DB: db}, nil
}

func (s *Store) Close() error {
	return s.DB.Close()
}

func isInvalidDatabaseValue(err error) bool {
	type sqlState interface{ SQLState() string }
	var state sqlState
	if !errors.As(err, &state) {
		return false
	}
	switch state.SQLState() {
	case "23503", "23514", "22P02", "22003":
		return true
	default:
		return false
	}
}
