package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Review struct {
	ID          int64  `json:"id"`
	PizzaID     int64  `json:"pizza_id"`
	PizzaSlug   string `json:"pizza_slug,omitempty"`
	UserID      int64  `json:"-"`
	DisplayName string `json:"display_name"`
	Rating      int    `json:"rating"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Status      string `json:"status,omitempty"`
	CreatedAt   string `json:"created_at"`
}

func (s *Store) ListReviews(ctx context.Context, slug string, page, pageSize int) ([]Review, int64, error) {
	var total int64
	err := s.DB.QueryRowContext(ctx, `
		SELECT COUNT(r.id)
		FROM pizzas p
		LEFT JOIN reviews r ON r.pizza_id = p.id AND r.status = 'published'
		WHERE p.slug = $1 AND p.status = 'published'
		GROUP BY p.id`, slug).Scan(&total)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, ErrNotFound
	}
	if err != nil {
		return nil, 0, fmt.Errorf("count reviews: %w", err)
	}
	rows, err := s.DB.QueryContext(ctx, `
		SELECT r.id, r.pizza_id, r.user_id, u.display_name, r.rating, r.title, r.body, r.created_at
		FROM reviews r
		JOIN users u ON u.id = r.user_id
		JOIN pizzas p ON p.id = r.pizza_id
		WHERE p.slug = $1 AND p.status = 'published' AND r.status = 'published'
		ORDER BY r.created_at DESC, r.id DESC
		LIMIT $2 OFFSET $3`, slug, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list reviews: %w", err)
	}
	defer rows.Close()
	reviews := make([]Review, 0)
	for rows.Next() {
		var review Review
		var createdAt sql.NullTime
		if err := rows.Scan(&review.ID, &review.PizzaID, &review.UserID, &review.DisplayName, &review.Rating, &review.Title, &review.Body, &createdAt); err != nil {
			return nil, 0, fmt.Errorf("scan review: %w", err)
		}
		if createdAt.Valid {
			review.CreatedAt = createdAt.Time.UTC().Format("2006-01-02T15:04:05Z")
		}
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate reviews: %w", err)
	}
	return reviews, total, nil
}

func (s *Store) SaveReview(ctx context.Context, slug string, userID int64, rating int, title, body string) (Review, error) {
	var review Review
	var createdAt sql.NullTime
	err := s.DB.QueryRowContext(ctx, `
		INSERT INTO reviews (pizza_id, user_id, rating, title, body, status)
		SELECT p.id, $2, $3, $4, $5, 'pending'
		FROM pizzas p
		WHERE p.slug = $1 AND p.status = 'published'
		ON CONFLICT (pizza_id, user_id)
		DO UPDATE SET rating = EXCLUDED.rating, title = EXCLUDED.title,
		              body = EXCLUDED.body, status = 'pending', updated_at = NOW()
		RETURNING id, pizza_id, user_id, rating, title, body, status, created_at`,
		slug, userID, rating, title, body,
	).Scan(&review.ID, &review.PizzaID, &review.UserID, &review.Rating, &review.Title, &review.Body, &review.Status, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Review{}, ErrNotFound
	}
	if err != nil {
		return Review{}, fmt.Errorf("save review: %w", err)
	}
	if createdAt.Valid {
		review.CreatedAt = createdAt.Time.UTC().Format("2006-01-02T15:04:05Z")
	}
	return review, nil
}
