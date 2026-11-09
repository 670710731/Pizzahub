package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func (s *Store) ListCategories(ctx context.Context) ([]Category, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id, name, slug, description FROM categories ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()
	categories := make([]Category, 0)
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Slug, &category.Description); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate categories: %w", err)
	}
	return categories, nil
}

func (s *Store) CreateCategory(ctx context.Context, category Category) (Category, error) {
	err := s.DB.QueryRowContext(ctx, `
		INSERT INTO categories (name, slug, description)
		VALUES ($1, $2, $3) RETURNING id`,
		strings.TrimSpace(category.Name), strings.TrimSpace(category.Slug), category.Description,
	).Scan(&category.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return Category{}, ErrConflict
		}
		return Category{}, fmt.Errorf("create category: %w", err)
	}
	return category, nil
}

func (s *Store) UpdateCategory(ctx context.Context, id int64, category Category) error {
	result, err := s.DB.ExecContext(ctx, `
		UPDATE categories SET name = $1, slug = $2, description = $3, updated_at = NOW()
		WHERE id = $4`,
		strings.TrimSpace(category.Name), strings.TrimSpace(category.Slug), category.Description, id,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrConflict
		}
		return fmt.Errorf("update category: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check updated category: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) DeleteCategory(ctx context.Context, id int64) error {
	result, err := s.DB.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check deleted category: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) ListAdminReviews(ctx context.Context, status string, page, pageSize int) ([]Review, int64, error) {
	if status != "pending" && status != "published" && status != "rejected" {
		return nil, 0, ErrInvalidInput
	}
	var total int64
	if err := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM reviews WHERE status = $1`, status).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count admin reviews: %w", err)
	}
	rows, err := s.DB.QueryContext(ctx, `
		SELECT r.id, r.pizza_id, p.slug, r.user_id, u.display_name,
		       r.rating, r.title, r.body, r.status, r.created_at
		FROM reviews r
		JOIN users u ON u.id = r.user_id
		JOIN pizzas p ON p.id = r.pizza_id
		WHERE r.status = $1
		ORDER BY r.created_at DESC, r.id DESC
		LIMIT $2 OFFSET $3`, status, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list admin reviews: %w", err)
	}
	defer rows.Close()
	reviews := make([]Review, 0)
	for rows.Next() {
		var review Review
		var createdAt sql.NullTime
		if err := rows.Scan(
			&review.ID, &review.PizzaID, &review.PizzaSlug, &review.UserID, &review.DisplayName,
			&review.Rating, &review.Title, &review.Body, &review.Status, &createdAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan admin review: %w", err)
		}
		if createdAt.Valid {
			review.CreatedAt = createdAt.Time.UTC().Format("2006-01-02T15:04:05Z")
		}
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate admin reviews: %w", err)
	}
	return reviews, total, nil
}

func (s *Store) ModerateReview(ctx context.Context, id int64, status string) error {
	result, err := s.DB.ExecContext(ctx, `
		UPDATE reviews SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, id)
	if err != nil {
		return fmt.Errorf("moderate review: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check moderated review: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
