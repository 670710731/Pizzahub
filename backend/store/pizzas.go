package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Pizza struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	Slug          string     `json:"slug"`
	Description   string     `json:"description"`
	Status        string     `json:"status"`
	RatingAverage float64    `json:"rating_average"`
	ReviewCount   int64      `json:"review_count"`
	Categories    []Category `json:"categories"`
	Images        []Image    `json:"images"`
	Variants      []Variant  `json:"variants"`
}

type Category struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type Image struct {
	URL       string `json:"image_url"`
	AltText   string `json:"alt_text"`
	SortOrder int    `json:"sort_order"`
}

type Variant struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name"`
	SizeLabel   string `json:"size_label"`
	Price       string `json:"price"`
	Currency    string `json:"currency"`
	IsAvailable bool   `json:"is_available"`
	SortOrder   int    `json:"sort_order"`
}

type VariantInput struct {
	Name        string `json:"name"`
	SizeLabel   string `json:"size_label"`
	Price       string `json:"price"`
	Currency    string `json:"currency"`
	IsAvailable *bool  `json:"is_available"`
	SortOrder   int    `json:"sort_order"`
}

type PizzaInput struct {
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	CategoryIDs []int64        `json:"category_ids"`
	Images      []Image        `json:"images"`
	Variants    []VariantInput `json:"variants"`
}

func (s *Store) ListPizzas(ctx context.Context, page, pageSize int, category, sort string) ([]Pizza, int64, error) {
	orderBy := "p.created_at DESC, p.id DESC"
	switch sort {
	case "name":
		orderBy = "p.name ASC, p.id ASC"
	case "rating":
		orderBy = "rating_average DESC, p.id DESC"
	case "newest", "":
	default:
		return nil, 0, ErrInvalidInput
	}
	var total int64
	err := s.DB.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT p.id)
		FROM pizzas p
		WHERE p.status = 'published'
		  AND ($1 = '' OR EXISTS (
			SELECT 1
			FROM pizza_categories pc
			JOIN categories c ON c.id = pc.category_id
			WHERE pc.pizza_id = p.id AND c.slug = $1
		  ))`, category).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count pizzas: %w", err)
	}
	query := fmt.Sprintf(`
		SELECT p.id, p.name, p.slug, p.description, p.status,
		       COALESCE(rating.rating_average, 0)::float8, COALESCE(rating.review_count, 0)
		FROM pizzas p
		LEFT JOIN LATERAL (
			SELECT AVG(r.rating) AS rating_average, COUNT(*) AS review_count
			FROM reviews r
			WHERE r.pizza_id = p.id AND r.status = 'published'
		) rating ON TRUE
		WHERE p.status = 'published'
		  AND ($1 = '' OR EXISTS (
			SELECT 1
			FROM pizza_categories pc
			JOIN categories c ON c.id = pc.category_id
			WHERE pc.pizza_id = p.id AND c.slug = $1
		  ))
		ORDER BY %s
		LIMIT $2 OFFSET $3`, orderBy)
	rows, err := s.DB.QueryContext(ctx, query, category, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list pizzas: %w", err)
	}
	defer rows.Close()
	pizzas := make([]Pizza, 0)
	for rows.Next() {
		var pizza Pizza
		if err := rows.Scan(&pizza.ID, &pizza.Name, &pizza.Slug, &pizza.Description, &pizza.Status, &pizza.RatingAverage, &pizza.ReviewCount); err != nil {
			return nil, 0, fmt.Errorf("scan pizza: %w", err)
		}
		pizzas = append(pizzas, pizza)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate pizzas: %w", err)
	}
	return pizzas, total, nil
}

func (s *Store) PizzaBySlug(ctx context.Context, slug string) (Pizza, error) {
	var pizza Pizza
	err := s.DB.QueryRowContext(ctx, `
		SELECT p.id, p.name, p.slug, p.description, p.status,
		       COALESCE(AVG(r.rating), 0)::float8, COUNT(r.id)
		FROM pizzas p
		LEFT JOIN reviews r ON r.pizza_id = p.id AND r.status = 'published'
		WHERE p.slug = $1 AND p.status = 'published'
		GROUP BY p.id`, slug).Scan(
		&pizza.ID, &pizza.Name, &pizza.Slug, &pizza.Description, &pizza.Status,
		&pizza.RatingAverage, &pizza.ReviewCount,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Pizza{}, ErrNotFound
	}
	if err != nil {
		return Pizza{}, fmt.Errorf("get pizza: %w", err)
	}
	pizza.Categories = make([]Category, 0)
	pizza.Images = make([]Image, 0)
	pizza.Variants = make([]Variant, 0)
	categoryRows, err := s.DB.QueryContext(ctx, `
		SELECT c.id, c.name, c.slug, c.description
		FROM categories c JOIN pizza_categories pc ON pc.category_id = c.id
		WHERE pc.pizza_id = $1 ORDER BY c.name`, pizza.ID)
	if err != nil {
		return Pizza{}, fmt.Errorf("get pizza categories: %w", err)
	}
	for categoryRows.Next() {
		var category Category
		if err := categoryRows.Scan(&category.ID, &category.Name, &category.Slug, &category.Description); err != nil {
			categoryRows.Close()
			return Pizza{}, fmt.Errorf("scan pizza category: %w", err)
		}
		pizza.Categories = append(pizza.Categories, category)
	}
	if err := categoryRows.Err(); err != nil {
		categoryRows.Close()
		return Pizza{}, fmt.Errorf("iterate pizza categories: %w", err)
	}
	categoryRows.Close()
	imageRows, err := s.DB.QueryContext(ctx, `
		SELECT image_url, alt_text, sort_order FROM pizza_images
		WHERE pizza_id = $1 ORDER BY sort_order`, pizza.ID)
	if err != nil {
		return Pizza{}, fmt.Errorf("get pizza images: %w", err)
	}
	for imageRows.Next() {
		var image Image
		if err := imageRows.Scan(&image.URL, &image.AltText, &image.SortOrder); err != nil {
			imageRows.Close()
			return Pizza{}, fmt.Errorf("scan pizza image: %w", err)
		}
		pizza.Images = append(pizza.Images, image)
	}
	if err := imageRows.Err(); err != nil {
		imageRows.Close()
		return Pizza{}, fmt.Errorf("iterate pizza images: %w", err)
	}
	imageRows.Close()
	variantRows, err := s.DB.QueryContext(ctx, `
		SELECT id, name, size_label, price::text, currency, is_available, sort_order
		FROM pizza_variants WHERE pizza_id = $1 AND is_available = TRUE
		ORDER BY sort_order`, pizza.ID)
	if err != nil {
		return Pizza{}, fmt.Errorf("get pizza variants: %w", err)
	}
	defer variantRows.Close()
	for variantRows.Next() {
		var variant Variant
		if err := variantRows.Scan(&variant.ID, &variant.Name, &variant.SizeLabel, &variant.Price, &variant.Currency, &variant.IsAvailable, &variant.SortOrder); err != nil {
			return Pizza{}, fmt.Errorf("scan pizza variant: %w", err)
		}
		pizza.Variants = append(pizza.Variants, variant)
	}
	if err := variantRows.Err(); err != nil {
		return Pizza{}, fmt.Errorf("iterate pizza variants: %w", err)
	}
	return pizza, nil
}

func (s *Store) SavePizza(ctx context.Context, id int64, input PizzaInput) (int64, error) {
	if input.Status == "" {
		input.Status = "draft"
	}
	if input.Status != "draft" && input.Status != "published" && input.Status != "archived" {
		return 0, ErrInvalidInput
	}
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin pizza transaction: %w", err)
	}
	defer tx.Rollback()
	if id == 0 {
		err = tx.QueryRowContext(ctx, `
			INSERT INTO pizzas (name, slug, description, status)
			VALUES ($1, $2, $3, $4) RETURNING id`,
			strings.TrimSpace(input.Name), strings.TrimSpace(input.Slug), input.Description, input.Status,
		).Scan(&id)
	} else {
		result, updateErr := tx.ExecContext(ctx, `
			UPDATE pizzas SET name = $1, slug = $2, description = $3, status = $4, updated_at = NOW()
			WHERE id = $5`,
			strings.TrimSpace(input.Name), strings.TrimSpace(input.Slug), input.Description, input.Status, id,
		)
		err = updateErr
		if err == nil {
			n, affectedErr := result.RowsAffected()
			if affectedErr != nil {
				return 0, fmt.Errorf("check updated pizza: %w", affectedErr)
			}
			if n == 0 {
				return 0, ErrNotFound
			}
		}
	}
	if err != nil {
		if isUniqueViolation(err) {
			return 0, ErrConflict
		}
		if isInvalidDatabaseValue(err) {
			return 0, ErrInvalidInput
		}
		return 0, fmt.Errorf("save pizza: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM pizza_categories WHERE pizza_id = $1`, id); err != nil {
		return 0, fmt.Errorf("replace pizza categories: %w", err)
	}
	for _, categoryID := range input.CategoryIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO pizza_categories (pizza_id, category_id) VALUES ($1, $2)`, id, categoryID); err != nil {
			if isInvalidDatabaseValue(err) {
				return 0, ErrInvalidInput
			}
			return 0, fmt.Errorf("assign pizza category: %w", err)
		}
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM pizza_images WHERE pizza_id = $1`, id); err != nil {
		return 0, fmt.Errorf("replace pizza images: %w", err)
	}
	for i, image := range input.Images {
		if image.SortOrder == 0 && i > 0 {
			image.SortOrder = i
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO pizza_images (pizza_id, image_url, alt_text, sort_order)
			VALUES ($1, $2, $3, $4)`, id, image.URL, image.AltText, image.SortOrder); err != nil {
			if isUniqueViolation(err) {
				return 0, ErrConflict
			}
			if isInvalidDatabaseValue(err) {
				return 0, ErrInvalidInput
			}
			return 0, fmt.Errorf("save pizza image: %w", err)
		}
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM pizza_variants WHERE pizza_id = $1`, id); err != nil {
		return 0, fmt.Errorf("replace pizza variants: %w", err)
	}
	for i, variant := range input.Variants {
		if variant.Currency == "" {
			variant.Currency = "THB"
		}
		if variant.SortOrder == 0 && i > 0 {
			variant.SortOrder = i
		}
		isAvailable := true
		if variant.IsAvailable != nil {
			isAvailable = *variant.IsAvailable
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO pizza_variants (pizza_id, name, size_label, price, currency, is_available, sort_order)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			id, variant.Name, variant.SizeLabel, variant.Price, variant.Currency, isAvailable, variant.SortOrder,
		); err != nil {
			if isUniqueViolation(err) {
				return 0, ErrConflict
			}
			if isInvalidDatabaseValue(err) {
				return 0, ErrInvalidInput
			}
			return 0, fmt.Errorf("save pizza variant: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		if isUniqueViolation(err) {
			return 0, ErrConflict
		}
		return 0, fmt.Errorf("commit pizza: %w", err)
	}
	return id, nil
}

func (s *Store) DeletePizza(ctx context.Context, id int64) error {
	result, err := s.DB.ExecContext(ctx, `
		UPDATE pizzas SET status = 'archived', updated_at = NOW()
		WHERE id = $1 AND status <> 'archived'`, id)
	if err != nil {
		return fmt.Errorf("delete pizza: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check deleted pizza: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
