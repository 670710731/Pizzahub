package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
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
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse database configuration: %w", err)
	}
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	db := stdlib.OpenDB(*config)
	maxOpenConns := 5
	if configured := os.Getenv("DB_MAX_OPEN_CONNS"); configured != "" {
		maxOpenConns, err = strconv.Atoi(configured)
		if err != nil || maxOpenConns < 1 {
			_ = db.Close()
			return nil, errors.New("DB_MAX_OPEN_CONNS must be a positive integer")
		}
	}
	db.SetMaxOpenConns(maxOpenConns)
	if maxOpenConns < 5 {
		db.SetMaxIdleConns(maxOpenConns)
	} else {
		db.SetMaxIdleConns(2)
	}
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
