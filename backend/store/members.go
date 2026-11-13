package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	DisplayName  string `json:"display_name"`
	Role         string `json:"role"`
	PasswordHash string `json:"-"`
}

func (s *Store) CreateUser(ctx context.Context, email, passwordHash, displayName string) (User, error) {
	var user User
	err := s.DB.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, display_name)
		VALUES (LOWER($1), $2, $3)
		RETURNING id, email, display_name, role`,
		strings.TrimSpace(email), passwordHash, strings.TrimSpace(displayName),
	).Scan(&user.ID, &user.Email, &user.DisplayName, &user.Role)
	if err != nil {
		if isUniqueViolation(err) {
			return User{}, ErrConflict
		}
		return User{}, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (s *Store) UserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := s.DB.QueryRowContext(ctx, `
		SELECT id, email, display_name, role, password_hash
		FROM users
		WHERE LOWER(email) = LOWER($1) AND status = 'active'`,
		strings.TrimSpace(email),
	).Scan(&user.ID, &user.Email, &user.DisplayName, &user.Role, &user.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("find user: %w", err)
	}
	return user, nil
}

func (s *Store) UserByID(ctx context.Context, id int64) (User, error) {
	var user User
	err := s.DB.QueryRowContext(ctx, `
		SELECT id, email, display_name, role, password_hash
		FROM users
		WHERE id = $1 AND status = 'active'`,
		id,
	).Scan(&user.ID, &user.Email, &user.DisplayName, &user.Role, &user.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("find user by id: %w", err)
	}
	return user, nil
}

func (s *Store) UpdatePasswordHash(ctx context.Context, id int64, currentHash, newHash string) error {
	result, err := s.DB.ExecContext(ctx, `
		UPDATE users SET password_hash = $1, updated_at = NOW()
		WHERE id = $2 AND password_hash = $3 AND status = 'active'`,
		newHash, id, currentHash,
	)
	if err != nil {
		return fmt.Errorf("update password hash: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check updated password: %w", err)
	}
	if rows == 0 {
		return ErrConflict
	}
	return nil
}

func isUniqueViolation(err error) bool {
	type sqlState interface{ SQLState() string }
	var state sqlState
	return errors.As(err, &state) && state.SQLState() == "23505"
}
