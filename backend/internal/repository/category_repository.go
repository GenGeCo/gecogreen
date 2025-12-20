package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gecogreen/backend/internal/models"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

// List returns all categories ordered by sort_order
func (r *CategoryRepository) List(ctx context.Context) ([]models.Category, error) {
	query := `
		SELECT id, name, slug, COALESCE(icon, ''), parent_id, sort_order
		FROM categories
		WHERE parent_id IS NULL
		ORDER BY sort_order, name
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Slug, &cat.Icon, &cat.ParentID, &cat.SortOrder)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

// GetByID returns a category by ID
func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	query := `
		SELECT id, name, slug, COALESCE(icon, ''), parent_id, sort_order
		FROM categories
		WHERE id = $1
	`

	var cat models.Category
	err := r.pool.QueryRow(ctx, query, id).Scan(&cat.ID, &cat.Name, &cat.Slug, &cat.Icon, &cat.ParentID, &cat.SortOrder)
	if err != nil {
		return nil, err
	}

	return &cat, nil
}

// GetBySlug returns a category by slug
func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	query := `
		SELECT id, name, slug, COALESCE(icon, ''), parent_id, sort_order
		FROM categories
		WHERE slug = $1
	`

	var cat models.Category
	err := r.pool.QueryRow(ctx, query, slug).Scan(&cat.ID, &cat.Name, &cat.Slug, &cat.Icon, &cat.ParentID, &cat.SortOrder)
	if err != nil {
		return nil, err
	}

	return &cat, nil
}
