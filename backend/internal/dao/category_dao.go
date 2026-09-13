package dao

import (
	"context"
	"database/sql"

	"agora-backend/internal/model"
)

type CategoryDAO struct {
	db *sql.DB
}

func NewCategoryDAO(db *sql.DB) *CategoryDAO {
	return &CategoryDAO{db: db}
}

func (d *CategoryDAO) ListCategories(ctx context.Context) ([]*model.Category, error) {
	query := `SELECT id, name, slug, description, sort_order, is_active, requires_review, created_at, updated_at FROM categories WHERE is_active = TRUE ORDER BY sort_order ASC, id ASC`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*model.Category, 0)
	for rows.Next() {
		c := &model.Category{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.SortOrder, &c.IsActive, &c.RequiresReview, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}
