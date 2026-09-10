package dao

import (
	"context"
	"database/sql"
	"errors"

	"Agora-BBS/internal/model"
)

type CategoryDAO struct {
	db *sql.DB
}

func NewCategoryDAO(db *sql.DB) *CategoryDAO {
	return &CategoryDAO{db: db}
}

func (d *CategoryDAO) Create(ctx context.Context, c *model.Category) error {
	query := `
		INSERT INTO categories (name, slug, description, parent_id, sort_order)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return d.db.QueryRowContext(ctx, query, c.Name, c.Slug, c.Description, c.ParentID, c.SortOrder).
		Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (d *CategoryDAO) ListAll(ctx context.Context) ([]*model.Category, error) {
	query := `
		SELECT id, name, slug, description, parent_id, sort_order, created_at, updated_at
		FROM categories
		ORDER BY sort_order ASC, id ASC
	`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Category
	for rows.Next() {
		c := &model.Category{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.ParentID, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (d *CategoryDAO) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	query := `
		SELECT id, name, slug, description, parent_id, sort_order, created_at, updated_at
		FROM categories WHERE id = $1
	`
	c := &model.Category{}
	err := d.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.ParentID, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return c, err
}
